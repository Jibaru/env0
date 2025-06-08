package scripts

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Jibaru/env0/pkg/client"
	"github.com/Jibaru/env0/pkg/envdiff"
	"github.com/Jibaru/env0/pkg/envfile"
	"github.com/Jibaru/env0/pkg/logger"
)

// PullInput represents the input parameters for the pull operation
type PullInput struct {
	TargetEnv *string
}

// PullFn represents a function that performs the pull operation
type PullFn func(context.Context, PullInput) error

// NewPull creates a new pull function with injected dependencies
func NewPull(c client.Client, logger logger.Logger) PullFn {
	return func(ctx context.Context, input PullInput) error {
		cfg, err := readConfigFile()
		if err != nil {
			return err
		}

		fullAppName := fmt.Sprintf("%s/%s", cfg.OwnerName, cfg.AppName)
		logger.Printf("pulling environments from app %s", fullAppName)

		envs, err := c.GetApp(ctx, fullAppName)
		if err != nil {
			return fmt.Errorf("failed to fetch environments: %v", err)
		}

		if err := processEnvironmentUpdates(envs, input.TargetEnv, logger); err != nil {
			return err
		}

		logger.Printf("environments pulled successfully")
		return nil
	}
}

func getEnvFileName(envName string) string {
	if envName == "" {
		return ".env"
	}
	return fmt.Sprintf(".env.%s", envName)
}

func processEnvironmentUpdates(envs map[string]map[string]interface{}, targetEnv *string, logger logger.Logger) error {
	for envName, remoteVars := range envs {
		if targetEnv != nil && envName != *targetEnv {
			continue
		}

		fileName := getEnvFileName(envName)

		// Load current environment if it exists
		currentVars, err := envfile.ParseEnvFile(fileName)
		if err != nil {
			if !os.IsNotExist(err) {
				return fmt.Errorf("failed to parse current env file %s: %v", fileName, err)
			}
			currentVars = make(map[string]interface{})
		}

		// Compare current and remote states
		diff := envdiff.CompareMaps(currentVars, remoteVars)

		if len(diff.Changes) == 0 {
			logger.Printf("no changes detected for environment: %s", envName)
			continue
		}

		// If there are changes, decide what to do
		if diff.SafeToMerge {
			// Only new variables, safe to merge
			mergedVars := envdiff.MergeMaps(currentVars, remoteVars, diff)
			if err := envfile.WriteEnvFile(fileName, mergedVars); err != nil {
				return fmt.Errorf("failed to write merged env file %s: %v", fileName, err)
			}
			logger.Printf("safely merged %d new variables into %s", len(diff.Changes), fileName)
		} else {
			// Write conflicts directly to the env file
			var content strings.Builder

			// First write all non-conflicting variables
			for k, v := range currentVars {
				isConflicting := false
				for _, change := range diff.Changes {
					if change.Name == k && change.Type == envdiff.Modified {
						isConflicting = true
						break
					}
				}
				if !isConflicting {
					content.WriteString(fmt.Sprintf("%s=%v\n", k, v))
				}
			}

			// Then write conflicts with git-style markers
			for _, change := range diff.Changes {
				if change.Type == envdiff.Modified {
					content.WriteString(envdiff.FormatGitStyleConflict(change.Name, change.OldValue, change.NewValue))
				} else if change.Type == envdiff.Added {
					content.WriteString(fmt.Sprintf("%s=%v\n", change.Name, change.NewValue))
				} else if change.Type == envdiff.Deleted {
					content.WriteString(envdiff.FormatGitStyleConflict(change.Name, change.OldValue, envdiff.DeletedValue{}))
				}
			}

			if err := os.WriteFile(fileName, []byte(content.String()), 0644); err != nil {
				return fmt.Errorf("failed to write conflict markers to %s: %v", fileName, err)
			}

			logger.Printf("detected conflicts in %s, marked conflicts in file with git-style markers", fileName)
			logger.Printf("please resolve conflicts manually and run push when ready")
		}
	}

	return nil
}

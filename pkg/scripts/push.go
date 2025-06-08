package scripts

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Jibaru/env0/pkg/client"
	"github.com/Jibaru/env0/pkg/envdiff"
	"github.com/Jibaru/env0/pkg/envfile"
	"github.com/Jibaru/env0/pkg/logger"
)

type config struct {
	AppName   string `json:"appName"`
	OwnerName string `json:"ownerName"`
}

// PushInput represents the input parameters for the push operation
type PushInput struct {
	TargetEnv *string
}

// PushFn represents a function that performs the push operation
type PushFn func(context.Context, PushInput) error

// NewPush creates a new push function with injected dependencies
func NewPush(c client.Client, logger logger.Logger) PushFn {
	return func(ctx context.Context, input PushInput) error {
		cfg, err := readConfigFile()
		if err != nil {
			return err
		}

		fullAppName := fmt.Sprintf("%s/%s", cfg.OwnerName, cfg.AppName)
		logger.Printf("reading environment files for app %s", fullAppName)

		// Get current remote state first
		remoteEnvs, err := c.GetApp(ctx, fullAppName)
		if err != nil {
			return fmt.Errorf("failed to fetch current remote state: %v", err)
		}

		// Process local environment files
		localEnvs, err := processEnvFiles(input.TargetEnv, logger)
		if err != nil {
			return err
		}

		// Compare and merge changes
		mergedEnvs, err := processPushUpdates(localEnvs, remoteEnvs, input.TargetEnv, logger)
		if err != nil {
			return err
		}

		if mergedEnvs != nil {
			logger.Printf("pushing environments to app %s", fullAppName)
			if err := c.UpdateApp(ctx, fullAppName, mergedEnvs); err != nil {
				return fmt.Errorf("failed to update environments: %v", err)
			}
			logger.Printf("environments pushed successfully")
		} else {
			logger.Printf("no changes to push")
		}

		return nil
	}
}

func processPushUpdates(localEnvs, remoteEnvs map[string]map[string]interface{}, targetEnv *string, logger logger.Logger) (map[string]map[string]interface{}, error) {
	mergedEnvs := make(map[string]map[string]interface{})
	hasChanges := false

	// Process each local environment
	for envName, localVars := range localEnvs {
		if targetEnv != nil && envName != *targetEnv {
			continue
		}

		remoteVars := remoteEnvs[envName]
		if remoteVars == nil {
			remoteVars = make(map[string]interface{})
		}

		// Compare local and remote states
		diff := envdiff.CompareMaps(remoteVars, localVars)

		if len(diff.Changes) == 0 {
			logger.Printf("no changes detected for environment: %s", envName)
			continue
		}

		// If there are changes, decide what to do
		if diff.SafeToMerge {
			// Only new variables, safe to merge
			mergedVars := envdiff.MergeMaps(remoteVars, localVars, diff)
			mergedEnvs[envName] = mergedVars
			hasChanges = true
			logger.Printf("safely merged %d new variables for environment: %s", len(diff.Changes), envName)
		} else {
			// There are conflicts, create conflict report
			report := envdiff.ConflictReport{
				Environment: envName,
				Conflicts:   diff.Changes,
				Timestamp:   time.Now().Format(time.RFC3339),
			}
			if err := envdiff.SaveConflictReport(report); err != nil {
				return nil, fmt.Errorf("failed to save conflict report for %s: %v", envName, err)
			}
			logger.Printf("detected conflicts in %s, saved conflict report", envName)
			logger.Printf("review conflicts in .env0/conflicts directory")
		}
	}

	if !hasChanges {
		return nil, nil
	}

	// Copy over any environments we didn't process
	for envName, vars := range remoteEnvs {
		if _, exists := mergedEnvs[envName]; !exists {
			mergedEnvs[envName] = vars
		}
	}

	return mergedEnvs, nil
}

func readConfigFile() (*config, error) {
	cfgData, err := os.ReadFile(filepath.Join(".env0", "config.json"))
	if err != nil {
		return nil, fmt.Errorf("app not initialized")
	}

	var cfg config
	if err := json.Unmarshal(cfgData, &cfg); err != nil {
		return nil, fmt.Errorf("invalid config file: %v", err)
	}

	return &cfg, nil
}

func getEnvNameFromFile(filename string) string {
	if filename == ".env" {
		return ""
	}
	return strings.TrimPrefix(filename, ".env.")
}

func processEnvFiles(targetEnv *string, logger logger.Logger) (map[string]map[string]interface{}, error) {
	envs := make(map[string]map[string]interface{})

	files, err := os.ReadDir(".")
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %v", err)
	}

	for _, fi := range files {
		name := fi.Name()
		if name != ".env" && !strings.HasPrefix(name, ".env.") {
			continue
		}

		envName := getEnvNameFromFile(name)
		if targetEnv != nil && envName != *targetEnv {
			continue
		}

		vars, err := envfile.ParseEnvFile(name)
		if err != nil {
			return nil, err
		}

		envs[envName] = vars
		logger.Printf("processed environment file: %s", name)
	}

	// If target set but not found, ensure empty
	if targetEnv != nil && envs[*targetEnv] == nil {
		envs[*targetEnv] = map[string]interface{}{}
		logger.Printf("creating empty environment for target: %s", *targetEnv)
	}

	return envs, nil
}

package scripts

import (
	"context"
	"fmt"
	"maps"
	"os"
	"slices"

	"github.com/Jibaru/env0/pkg/client"
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

		if err := writeEnvironmentFiles(envs, input.TargetEnv, logger); err != nil {
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

func writeEnvFile(fileName string, vars map[string]interface{}) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create environment file %s: %v", fileName, err)
	}
	defer file.Close()

	// Sort keys for consistent output
	keys := slices.Collect(maps.Keys(vars))
	slices.Sort(keys)

	// Write sorted key-value pairs
	for _, k := range keys {
		fmt.Fprintf(file, "%s=%v\n", k, vars[k])
	}

	return nil
}

func writeEnvironmentFiles(envs map[string]map[string]interface{}, targetEnv *string, logger logger.Logger) error {
	for envName, vars := range envs {
		if targetEnv != nil && envName != *targetEnv {
			continue
		}

		fileName := getEnvFileName(envName)
		if err := writeEnvFile(fileName, vars); err != nil {
			return err
		}

		logger.Printf("created environment file: %s", fileName)
	}

	return nil
}

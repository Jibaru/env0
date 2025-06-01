package scripts

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
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
		cfgData, err := os.ReadFile(filepath.Join(".env0", "config.json"))
		if err != nil {
			return fmt.Errorf("app not initialized")
		}

		var cfg struct {
			AppName   string `json:"appName"`
			OwnerName string `json:"ownerName"`
		}
		if err := json.Unmarshal(cfgData, &cfg); err != nil {
			return fmt.Errorf("invalid config file: %v", err)
		}

		fullAppName := fmt.Sprintf("%s/%s", cfg.OwnerName, cfg.AppName)

		logger.Printf("pulling environments from app %s", fullAppName)

		envs, err := c.GetApp(ctx, fullAppName)
		if err != nil {
			return fmt.Errorf("failed to fetch environments: %v", err)
		}

		// Write .env files
		for envName, vars := range envs {
			if input.TargetEnv != nil {
				if envName != *input.TargetEnv {
					continue
				}
			}

			fileName := ""
			if envName == "" {
				fileName = ".env"
			} else {
				fileName = fmt.Sprintf(".env.%s", envName)
			}

			file, err := os.Create(fileName)
			if err != nil {
				return fmt.Errorf("failed to create environment file %s: %v", fileName, err)
			}

			keys := slices.Collect(maps.Keys(vars))
			slices.Sort(keys)

			for k, v := range vars {
				fmt.Fprintf(file, "%s=%v\n", k, v)
			}
			file.Close()

			logger.Printf("created environment file: %s", fileName)
		}

		logger.Printf("environments pulled successfully")
		return nil
	}
}

package scripts

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Jibaru/env0/pkg/client"
	"github.com/Jibaru/env0/pkg/logger"
)

// PushInput represents the input parameters for the push operation
type PushInput struct {
	TargetEnv *string
}

// PushFn represents a function that performs the push operation
type PushFn func(context.Context, PushInput) error

// NewPush creates a new push function with injected dependencies
func NewPush(c client.Client, logger logger.Logger) PushFn {
	return func(ctx context.Context, input PushInput) error {
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

		logger.Printf("reading environment files for app %s", fullAppName)

		// Read .env files
		envs := make(map[string]map[string]interface{})
		files, err := os.ReadDir(".")
		if err != nil {
			return fmt.Errorf("failed to read directory: %v", err)
		}

		for _, fi := range files {
			name := fi.Name()
			if name != ".env" && !strings.HasPrefix(name, ".env.") {
				continue
			}

			var envName string
			if name == ".env" {
				envName = ""
			} else {
				envName = strings.TrimPrefix(name, ".env.")
			}

			if input.TargetEnv != nil && envName != *input.TargetEnv {
				continue
			}

			content, err := os.ReadFile(name)
			if err != nil {
				return fmt.Errorf("failed to read environment file %s: %v", name, err)
			}

			lines := strings.Split(string(content), "\n")
			vars := make(map[string]interface{})
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" {
					continue
				}
				parts := strings.SplitN(line, "=", 2)
				if len(parts) != 2 {
					logger.Printf("warning: invalid line format in %s: %s", name, line)
					continue
				}
				vars[parts[0]] = parts[1]
			}
			envs[envName] = vars
			logger.Printf("processed environment file: %s", name)
		}

		// If target set but not found, ensure empty
		if input.TargetEnv != nil && envs[*input.TargetEnv] == nil {
			envs[*input.TargetEnv] = map[string]interface{}{}
			logger.Printf("creating empty environment for target: %s", *input.TargetEnv)
		}

		logger.Printf("pushing environments to app %s", fullAppName)

		if err := c.UpdateApp(ctx, fullAppName, envs); err != nil {
			return fmt.Errorf("failed to update environments: %v", err)
		}

		logger.Printf("environments pushed successfully")
		return nil
	}
}

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

		envs, err := processEnvFiles(input.TargetEnv, logger)
		if err != nil {
			return err
		}

		logger.Printf("pushing environments to app %s", fullAppName)
		if err := c.UpdateApp(ctx, fullAppName, envs); err != nil {
			return fmt.Errorf("failed to update environments: %v", err)
		}

		logger.Printf("environments pushed successfully")
		return nil
	}
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

func parseEnvFile(filename string) (map[string]interface{}, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read environment file %s: %v", filename, err)
	}

	vars := make(map[string]interface{})
	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		vars[parts[0]] = parts[1]
	}

	return vars, nil
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

		vars, err := parseEnvFile(name)
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

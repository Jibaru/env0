package scripts

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Jibaru/env0/pkg/client"
	"github.com/Jibaru/env0/pkg/logger"
)

// InitInput represents the input parameters for the init operation
type InitInput struct {
	AppName string
}

// InitConfig represents the configuration structure for the app
type InitConfig struct {
	AppName   string `json:"appName"`
	OwnerName string `json:"ownerName"`
}

// InitFn represents a function that performs the init operation
type InitFn func(context.Context, InitInput) error

// NewInit creates a new init function with injected dependencies
func NewInit(c client.Client, logger logger.Logger) InitFn {
	return func(ctx context.Context, input InitInput) error {
		// Check if .env0 already exists
		if _, err := os.Stat(".env0"); err == nil {
			return fmt.Errorf("this app already has a configuration, use clone command instead")
		}

		logger.Printf("creating new app: %s", input.AppName)

		// Create the app via API
		ownerName, err := c.CreateApp(ctx, input.AppName)
		if err != nil {
			return fmt.Errorf("failed to create app: %v", err)
		}

		// Write local config
		cfg := InitConfig{
			AppName:   input.AppName,
			OwnerName: ownerName,
		}
		data, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal config: %v", err)
		}

		if err := os.Mkdir(".env0", 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %v", err)
		}

		if err := os.WriteFile(filepath.Join(".env0", "config.json"), data, 0644); err != nil {
			return fmt.Errorf("failed to write config file: %v", err)
		}

		logger.Printf("app %s created successfully", input.AppName)
		return nil
	}
}

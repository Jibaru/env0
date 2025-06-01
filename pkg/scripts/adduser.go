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

// AddUserInput represents the input parameters for the add user operation
type AddUserInput struct {
	Username string
}

// AddUserFn represents a function that performs the add user operation
type AddUserFn func(context.Context, AddUserInput) error

// NewAddUser creates a new add user function with injected dependencies
func NewAddUser(c client.Client, logger logger.Logger) AddUserFn {
	return func(ctx context.Context, input AddUserInput) error {
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

		logger.Printf("adding user %s to app %s", input.Username, fullAppName)

		if err := c.AddUser(ctx, fullAppName, input.Username); err != nil {
			return fmt.Errorf("failed to add user: %v", err)
		}

		logger.Printf("user %s successfully added to app %s", input.Username, fullAppName)
		return nil
	}
}

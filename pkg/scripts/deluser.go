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

// DeleteUserInput represents the input parameters for the delete user operation
type DeleteUserInput struct {
	Username string
}

// DeleteUserFn represents a function that performs the delete user operation
type DeleteUserFn func(context.Context, DeleteUserInput) error

// NewDeleteUser creates a new delete user function with injected dependencies
func NewDeleteUser(c client.Client, logger logger.Logger) DeleteUserFn {
	return func(ctx context.Context, input DeleteUserInput) error {
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

		logger.Printf("removing user %s from app %s", input.Username, fullAppName)

		if err := c.RemoveUser(ctx, fullAppName, input.Username); err != nil {
			return fmt.Errorf("failed to remove user: %v", err)
		}

		logger.Printf("user %s successfully removed from app %s", input.Username, fullAppName)
		return nil
	}
}

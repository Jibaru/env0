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

// ListUsersInput represents the input parameters for the list users operation
type ListUsersInput struct {
	// Empty since we'll get the app name from config
}

// ListUsersFn represents a function that performs the list users operation
type ListUsersFn func(context.Context, ListUsersInput) error

// NewListUsers creates a new list users function with injected dependencies
func NewListUsers(c client.Client, logger logger.Logger) ListUsersFn {
	return func(ctx context.Context, input ListUsersInput) error {
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
		logger.Printf("listing users for app %s", fullAppName)

		users, err := c.ListAppUsers(ctx, fullAppName)
		if err != nil {
			return fmt.Errorf("failed to list users: %v", err)
		}

		// Print each user's information
		for _, user := range users {
			logger.Printf("User: %s", user.Username)
			logger.Printf("  ID: %s", user.ID)
			logger.Printf("  Email: %s", user.Email)
			if user.IsOwner {
				logger.Printf("  Role: Owner")
			} else {
				logger.Printf("  Role: Collaborator")
			}
			logger.Printf("---")
		}

		if len(users) == 0 {
			logger.Printf("no users found")
		}

		return nil
	}
}

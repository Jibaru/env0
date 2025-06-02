package scripts

import (
	"context"
	"fmt"

	"github.com/Jibaru/env0/pkg/client"
	"github.com/Jibaru/env0/pkg/logger"
)

// ListAppsInput represents the input parameters for the list apps operation
type ListAppsInput struct {
	// Empty since we don't need any input parameters for listing all apps
}

// ListAppsFn represents a function that performs the list apps operation
type ListAppsFn func(context.Context, ListAppsInput) error

// NewListApps creates a new list apps function with injected dependencies
func NewListApps(c client.Client, logger logger.Logger) ListAppsFn {
	return func(ctx context.Context, input ListAppsInput) error {
		logger.Printf("listing all apps")

		// Call with no pagination to get all apps
		apps, err := c.ListApps(ctx, 0, 0, "desc", "")
		if err != nil {
			return fmt.Errorf("failed to list apps: %v", err)
		}

		// Print each app's information
		for _, app := range apps {
			logger.Printf("App: %s", app.Name)
			logger.Printf("  ID: %s", app.ID)
			logger.Printf("  Created: %s", app.CreatedAt)
			logger.Printf("  Environment count: %d", len(app.Envs))
			logger.Printf("  Other users count: %d", len(app.OtherUsersAllowedIds))
			logger.Printf("---")
		}

		if len(apps) == 0 {
			logger.Printf("no apps found")
		}

		return nil
	}
}

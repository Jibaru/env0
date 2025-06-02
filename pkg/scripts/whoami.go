package scripts

import (
	"context"

	"github.com/Jibaru/env0/pkg/auth"
	"github.com/Jibaru/env0/pkg/client"
	"github.com/Jibaru/env0/pkg/logger"
)

// WhoAmIInput represents the input parameters for the whoami operation
type WhoAmIInput struct {
	// Empty since we don't need any input parameters
}

// WhoAmIFn represents a function that performs the whoami operation
type WhoAmIFn func(context.Context, WhoAmIInput) error

// NewWhoAmI creates a new whoami function with injected dependencies
func NewWhoAmI(c client.Client, logger logger.Logger) WhoAmIFn {
	return func(ctx context.Context, input WhoAmIInput) error {
		authData, err := auth.Load()
		if err != nil {
			logger.Printf("Status: Not authenticated")
			return nil
		}

		// If we have user info, display it
		if authData.User.Username != "" {
			logger.Printf("Username: %s", authData.User.Username)
			logger.Printf("Email: %s", authData.User.Email)
			logger.Printf("Status: Authenticated")
		} else {
			// We have a token but no user info, show limited info
			logger.Printf("Status: Authenticated")
			logger.Printf("Note: Login again to see full user information")
		}

		return nil
	}
}

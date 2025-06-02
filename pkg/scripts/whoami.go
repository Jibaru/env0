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
			logger.Printf("Reason: %v", err)
			return nil
		}

		if !authData.IsAuthenticated() {
			logger.Printf("Status: Not authenticated")
			logger.Printf("Reason: Token is invalid or expired")
			return nil
		}

		logger.Printf("Status: Authenticated")

		if authData.HasUserInfo() {
			logger.Printf("Username: %s", authData.User.Username)
			logger.Printf("Email: %s", authData.User.Email)
		} else {
			logger.Printf("Note: Login again to see full user information")
		}

		return nil
	}
}

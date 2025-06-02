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

// WhoAmIInput represents the input parameters for the whoami operation
type WhoAmIInput struct {
	// Empty since we don't need any input parameters
}

// WhoAmIFn represents a function that performs the whoami operation
type WhoAmIFn func(context.Context, WhoAmIInput) error

// NewWhoAmI creates a new whoami function with injected dependencies
func NewWhoAmI(c client.Client, logger logger.Logger) WhoAmIFn {
	return func(ctx context.Context, input WhoAmIInput) error {
		home, err := getHomeDir()
		if err != nil {
			logger.Printf("Status: Not authenticated")
			return nil
		}

		// Try to read auth file
		authPath := filepath.Join(home, ".env0_cfg", "auth.json")
		authData, err := os.ReadFile(authPath)
		if err != nil {
			logger.Printf("Status: Not authenticated")
			return nil
		}

		var auth struct {
			Token string `json:"token"`
			User  struct {
				Username string `json:"username"`
				Email    string `json:"email"`
			} `json:"user"`
		}

		if err := json.Unmarshal(authData, &auth); err != nil {
			logger.Printf("Status: Not authenticated")
			return nil
		}

		if auth.Token == "" {
			logger.Printf("Status: Not authenticated")
			return nil
		}

		// If we have user info, display it
		if auth.User.Username != "" {
			logger.Printf("Username: %s", auth.User.Username)
			logger.Printf("Email: %s", auth.User.Email)
			logger.Printf("Status: Authenticated")
		} else {
			// We have a token but no user info, show limited info
			logger.Printf("Status: Authenticated")
			logger.Printf("Note: Login again to see full user information")
		}

		return nil
	}
}

// getHomeDir returns the user's home directory in a cross-platform way
func getHomeDir() (string, error) {
	// Try USERPROFILE for Windows first
	if home := os.Getenv("USERPROFILE"); home != "" {
		return home, nil
	}
	// Try HOME for Unix-like systems
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}
	return "", fmt.Errorf("unable to determine user home directory")
}

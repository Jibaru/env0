package scripts

import (
	"context"

	"github.com/Jibaru/env0/pkg/client"
)

// LoginInput represents the input parameters for the login operation
type LoginInput struct {
	UsernameOrEmail string
	Password        string
}

// LoginFn represents a function that performs the login operation
type LoginFn func(context.Context, LoginInput) error

// NewLogin creates a new login function with injected dependencies
func NewLogin(c client.Client, logger Logger) LoginFn {
	return func(ctx context.Context, input LoginInput) error {
		logger.Printf("attempting login for user: %s", input.UsernameOrEmail)

		if err := c.Login(ctx, input.UsernameOrEmail, input.Password); err != nil {
			return err
		}

		logger.Printf("login successful")
		logger.Printf("token saved to configuration")

		return nil
	}
}

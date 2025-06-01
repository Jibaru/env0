package scripts

import (
	"context"
	"fmt"

	"github.com/Jibaru/env0/pkg/client"
	"github.com/Jibaru/env0/pkg/logger"
)

// SignupInput represents the input parameters for the signup operation
type SignupInput struct {
	Username string
	Email    string
	Password string
}

// SignupFn represents a function that performs the signup operation
type SignupFn func(context.Context, SignupInput) error

// NewSignup creates a new signup function with injected dependencies
func NewSignup(c client.Client, logger logger.Logger) SignupFn {
	return func(ctx context.Context, input SignupInput) error {
		logger.Printf("attempting to create account for user: %s", input.Username)

		if err := c.Signup(ctx, input.Username, input.Email, input.Password); err != nil {
			return fmt.Errorf("failed to create account: %v", err)
		}

		logger.Printf("account created successfully for user: %s", input.Username)
		return nil
	}
}

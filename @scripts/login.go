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

// LoginOutput represents the result of the login operation
type LoginOutput struct {
	Message string
	Error   error
}

// LoginDependencies contains the required dependencies for the login operation
type LoginDependencies struct {
	Client *client.Client
}

// Login handles the authentication process with Env0
func Login(ctx context.Context, deps LoginDependencies, input LoginInput) LoginOutput {
	if err := deps.Client.Login(ctx, input.UsernameOrEmail, input.Password); err != nil {
		return LoginOutput{
			Message: err.Error(),
			Error:   err,
		}
	}

	return LoginOutput{
		Message: "Login successful",
		Error:   nil,
	}
}

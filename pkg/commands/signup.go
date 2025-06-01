package commands

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/Jibaru/env0/pkg/scripts"
)

func signupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "signup <username> <email> <password>",
		Args:  cobra.ExactArgs(3),
		Short: "Create a new Env0 account",
		RunE: func(cmd *cobra.Command, args []string) error {
			signup := scripts.NewSignup(apiClient, logger)
			return signup(context.Background(), scripts.SignupInput{
				Username: args[0],
				Email:    args[1],
				Password: args[2],
			})
		},
	}
	return cmd
}

package commands

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/Jibaru/env0/pkg/scripts"
)

func loginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login <usernameOrEmail> <password>",
		Args:  cobra.ExactArgs(2),
		Short: "Authenticate with Env0",
		RunE: func(cmd *cobra.Command, args []string) error {
			login := scripts.NewLogin(apiClient, logger)
			return login(context.Background(), scripts.LoginInput{
				UsernameOrEmail: args[0],
				Password:        args[1],
			})
		},
	}
	return cmd
}

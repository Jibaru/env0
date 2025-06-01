package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Jibaru/env0/pkg/client"
	"github.com/Jibaru/env0/pkg/scripts"
)

func delUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deluser <username>",
		Args:  cobra.ExactArgs(1),
		Short: "Remove a user from the initialized Env0 app",
		RunE: func(cmd *cobra.Command, args []string) error {
			token, err := scripts.LoadAndValidateToken()
			if err != nil {
				return fmt.Errorf("authentication required")
			}

			authClient := client.New(token)

			deleteUser := scripts.NewDeleteUser(authClient, logger)
			return deleteUser(context.Background(), scripts.DeleteUserInput{
				Username: args[0],
			})
		},
	}
	return cmd
}

package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Jibaru/env0/pkg/client"
	"github.com/Jibaru/env0/pkg/scripts"
)

func addUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "adduser <username>",
		Args:  cobra.ExactArgs(1),
		Short: "Add a user to the initialized Env0 app",
		RunE: func(cmd *cobra.Command, args []string) error {
			token, err := scripts.LoadAndValidateToken()
			if err != nil {
				return fmt.Errorf("authentication required")
			}

			authClient := client.New(token)

			addUser := scripts.NewAddUser(authClient, logger)
			return addUser(context.Background(), scripts.AddUserInput{
				Username: args[0],
			})
		},
	}
	return cmd
}

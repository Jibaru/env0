package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Jibaru/env0/pkg/client"
	"github.com/Jibaru/env0/pkg/scripts"
)

func listUsersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listusers",
		Args:  cobra.NoArgs,
		Short: "List all users with access to the initialized Env0 app",
		RunE: func(cmd *cobra.Command, args []string) error {
			token, err := scripts.LoadAndValidateToken()
			if err != nil {
				return fmt.Errorf("authentication required")
			}

			authClient := client.New(token)

			listUsers := scripts.NewListUsers(authClient, logger)
			return listUsers(context.Background(), scripts.ListUsersInput{})
		},
	}
	return cmd
}

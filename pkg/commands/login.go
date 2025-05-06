package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"env0/pkg/client"
)

func loginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login <usernameOrEmail> <password>",
		Args:  cobra.ExactArgs(2),
		Short: "Authenticate with Env0",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			usernameOrEmail, password := args[0], args[1]
			c := client.New("") // token will be populated on successful login

			if err := c.Login(ctx, usernameOrEmail, password); err != nil {
				// Print API error message if available
				fmt.Println(err.Error())
				return nil
			}

			fmt.Println("Login ok")
			return nil
		},
	}
	return cmd
}

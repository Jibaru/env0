package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"env0/pkg/client"
)

func signupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "signup <username> <email> <password>",
		Args:  cobra.ExactArgs(3),
		Short: "Create a new Env0 account",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			c := client.New("") // unauthenticated client
			username, email, password := args[0], args[1], args[2]

			if err := c.Signup(ctx, username, email, password); err != nil {
				// If it's a ClientError, print its message; else, generic failure
				fmt.Println(err.Error())
				return nil
			}

			fmt.Println("Signup ok")
			return nil
		},
	}
	return cmd
}

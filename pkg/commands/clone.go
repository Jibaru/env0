package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Jibaru/env0/pkg/client"
	"github.com/Jibaru/env0/pkg/scripts"
)

func cloneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clone <fullAppName>",
		Args:  cobra.ExactArgs(1),
		Short: "Clone an existing Env0 app's environments",
		RunE: func(cmd *cobra.Command, args []string) error {
			token, err := scripts.LoadAndValidateToken()
			if err != nil {
				return fmt.Errorf("authentication required")
			}

			authClient := client.New(token)

			clone := scripts.NewClone(authClient, logger)
			return clone(context.Background(), scripts.CloneInput{
				FullAppName: args[0],
			})
		},
	}
	return cmd
}

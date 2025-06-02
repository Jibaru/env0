package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Jibaru/env0/pkg/client"
	"github.com/Jibaru/env0/pkg/scripts"
)

func listAppsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listapps",
		Args:  cobra.NoArgs,
		Short: "List all Env0 apps you have access to",
		RunE: func(cmd *cobra.Command, args []string) error {
			token, err := scripts.LoadAndValidateToken()
			if err != nil {
				return fmt.Errorf("authentication required")
			}

			authClient := client.New(token)

			listApps := scripts.NewListApps(authClient, logger)
			return listApps(context.Background(), scripts.ListAppsInput{})
		},
	}
	return cmd
}

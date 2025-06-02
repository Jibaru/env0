package commands

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/Jibaru/env0/pkg/scripts"
)

func whoamiCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whoami",
		Args:  cobra.NoArgs,
		Short: "Display information about the current user",
		RunE: func(cmd *cobra.Command, args []string) error {
			// We don't validate token here since we want to show "not authenticated" status
			whoami := scripts.NewWhoAmI(apiClient, logger)
			return whoami(context.Background(), scripts.WhoAmIInput{})
		},
	}
	return cmd
}

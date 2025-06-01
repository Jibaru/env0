package commands

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/Jibaru/env0/pkg/scripts"
)

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Args:  cobra.NoArgs,
		Short: "Get version",
		RunE: func(cmd *cobra.Command, args []string) error {
			version := scripts.NewVersion(logger)
			return version(context.Background())
		},
	}
}

package commands

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/Jibaru/env0/pkg/scripts"
)

func cfgCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cfg",
		Args:  cobra.NoArgs,
		Short: "Show the env0 configuration directory path and status",
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath := scripts.NewConfigPath(apiClient, logger)
			return configPath(context.Background(), scripts.ConfigPathInput{})
		},
	}
	return cmd
}

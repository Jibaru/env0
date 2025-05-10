package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func versionCmd() *cobra.Command {
	const version = "v0.0.3"

	return &cobra.Command{
		Use:   "version",
		Args:  cobra.NoArgs,
		Short: "Get version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(version)
			return nil
		},
	}
}

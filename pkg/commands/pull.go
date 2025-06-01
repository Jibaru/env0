package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Jibaru/env0/pkg/client"
	"github.com/Jibaru/env0/pkg/scripts"
)

const defaultTargetEnv string = "default"

var defaultTargetEnvKey string = ""

func pullCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pull [envName]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Pull the latest environments for the initialized app",
		RunE: func(cmd *cobra.Command, args []string) error {
			var target *string
			if len(args) == 1 {
				target = &args[0]
				if *target == defaultTargetEnv {
					target = &defaultTargetEnvKey
				}
			}

			token, err := scripts.LoadAndValidateToken()
			if err != nil {
				return fmt.Errorf("authentication required")
			}

			authClient := client.New(token)

			pull := scripts.NewPull(authClient, logger)
			return pull(context.Background(), scripts.PullInput{
				TargetEnv: target,
			})
		},
	}
	return cmd
}

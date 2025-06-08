package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Jibaru/env0/pkg/client"
	"github.com/Jibaru/env0/pkg/scripts"
)

func pushCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "push [envName]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Push local environment files to the remote app",
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
			reader := bufio.NewReader(os.Stdin)

			push := scripts.NewPush(authClient, logger, reader)
			return push(context.Background(), scripts.PushInput{
				TargetEnv: target,
			})
		},
	}
	return cmd
}

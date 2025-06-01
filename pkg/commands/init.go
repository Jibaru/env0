package commands

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/Jibaru/env0/pkg/client"
	"github.com/Jibaru/env0/pkg/scripts"
)

func initCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init <appname>",
		Args:  cobra.ExactArgs(1),
		Short: "Initialize a new Env0 app in this directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			token, err := scripts.LoadAndValidateToken()
			if err != nil {
				return err
			}

			authClient := client.New(token)

			init := scripts.NewInit(authClient, logger)
			return init(context.Background(), scripts.InitInput{
				AppName: args[0],
			})
		},
	}
	return cmd
}

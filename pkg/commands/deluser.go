package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/Jibaru/env0/pkg/client"
)

func delUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deluser <username>",
		Args:  cobra.ExactArgs(1),
		Short: "Remove a user from the initialized Env0 app",
		RunE: func(cmd *cobra.Command, args []string) error {
			username := args[0]

			// 1) Auth
			token, err := client.LoadToken()
			if err != nil {
				fmt.Println("Authenticate first")
				return nil
			}
			if token == "" {
				fmt.Println("Authenticate again")
				return nil
			}

			// 2) Config
			cfgPath := filepath.Join(".env0", "config.json")
			data, err := os.ReadFile(cfgPath)
			if err != nil {
				fmt.Println("App not initialized")
				return nil
			}
			var cfg struct {
				AppName   string `json:"appName"`
				OwnerName string `json:"ownerName"`
			}
			_ = json.Unmarshal(data, &cfg)
			fullAppName := fmt.Sprintf("%s/%s", cfg.OwnerName, cfg.AppName)

			// 3) API call
			c := client.New(token)
			if err := c.RemoveUser(context.Background(), fullAppName, username); err != nil {
				fmt.Println(err.Error())
				return nil
			}

			fmt.Println("User removed")
			return nil
		},
	}
	return cmd
}

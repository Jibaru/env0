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

func addUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "adduser <username>",
		Args:  cobra.ExactArgs(1),
		Short: "Add a user to the initialized Env0 app",
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
			cfgData, err := os.ReadFile(filepath.Join(".env0", "config.json"))
			if err != nil {
				fmt.Println("App not initialized")
				return nil
			}
			var cfg struct {
				AppName   string `json:"appName"`
				OwnerName string `json:"ownerName"`
			}
			_ = json.Unmarshal(cfgData, &cfg)
			fullAppName := fmt.Sprintf("%s/%s", cfg.OwnerName, cfg.AppName)

			// 3) API call
			c := client.New(token)
			if err := c.AddUser(context.Background(), fullAppName, username); err != nil {
				// Print any API error
				fmt.Println(err.Error())
				return nil
			}

			fmt.Println("User added")
			return nil
		},
	}
	return cmd
}

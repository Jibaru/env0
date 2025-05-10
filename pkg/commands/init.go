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

type initConfig struct {
	AppName   string `json:"appName"`
	OwnerName string `json:"ownerName"`
}

func initCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init <appname>",
		Args:  cobra.ExactArgs(1),
		Short: "Initialize a new Env0 app in this directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			appName := args[0]

			// 1) Load auth token
			token, err := client.LoadToken()
			if err != nil {
				fmt.Println("Authenticate first")
				return nil
			}
			if token == "" {
				fmt.Println("Authenticate again")
				return nil
			}

			// 2) Check if .env0 already exists
			if _, err := os.Stat(".env0"); err == nil {
				fmt.Println("This app already has a configuration, use clone command instead")
				return nil
			}

			// 3) Create the app via API
			c := client.New(token)
			ownerName, err := c.CreateApp(context.Background(), appName)
			if err != nil {
				// ClientError will include status and message if set
				if ce, ok := err.(*client.ClientError); ok && ce.Err != nil {
					fmt.Println(ce.Err.Error())
				} else {
					fmt.Println("App creation failed")
				}
				return nil
			}

			// 4) On success, write local config
			cfg := initConfig{AppName: appName, OwnerName: ownerName}
			data, _ := json.MarshalIndent(cfg, "", "  ")

			if err := os.Mkdir(".env0", 0755); err != nil {
				return err
			}
			if err := os.WriteFile(filepath.Join(".env0", "config.json"), data, 0644); err != nil {
				return err
			}

			fmt.Println("App created")
			return nil
		},
	}
	return cmd
}

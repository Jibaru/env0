package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Jibaru/env0/pkg/client"
)

func cloneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clone <fullAppName>",
		Args:  cobra.ExactArgs(1),
		Short: "Clone an existing Env0 app's environments",
		RunE: func(cmd *cobra.Command, args []string) error {
			fullAppName := args[0]

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

			// 2) Check for existing config
			if _, err := os.Stat(".env0/config.json"); err == nil {
				fmt.Println("There is an app already cloned")
				return nil
			}

			// 3) Fetch envs from API
			c := client.New(token)
			envs, err := c.GetApp(context.Background(), fullAppName)
			if err != nil {
				// print API error message if ClientError
				fmt.Println(err.Error())
				return nil
			}

			// 4) Write .env files
			for envName, vars := range envs {
				fileName := fmt.Sprintf(".env.%s", envName)
				file, err := os.Create(fileName)
				if err != nil {
					return err
				}
				for k, v := range vars {
					fmt.Fprintf(file, "%s=%v\n", k, v)
				}
				file.Close()
			}

			// 5) Save local config
			parts := strings.SplitN(fullAppName, "/", 2)
			owner := parts[0]
			app := parts[1]
			if err := os.MkdirAll(".env0", 0755); err != nil {
				return err
			}
			cfgContent := fmt.Sprintf("{\"appName\": \"%s\", \"ownerName\": \"%s\"}", app, owner)
			os.WriteFile(filepath.Join(".env0", "config.json"), []byte(cfgContent), 0644)

			fmt.Println("App cloned")
			return nil
		},
	}
	return cmd
}

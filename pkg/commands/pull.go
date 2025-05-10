package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"

	"github.com/spf13/cobra"

	"github.com/Jibaru/env0/pkg/client"
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

			// 2) Check local config
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
			json.Unmarshal(data, &cfg)
			fullAppName := fmt.Sprintf("%s/%s", cfg.OwnerName, cfg.AppName)

			// 3) Fetch envs from API
			c := client.New(token)
			envs, err := c.GetApp(context.Background(), fullAppName)
			if err != nil {
				fmt.Println(err.Error())
				return nil
			}

			// 4) Write .env files
			for envName, vars := range envs {
				if target != nil {
					if envName != *target {
						continue
					}
				}

				fileName := ""
				if envName == "" {
					fileName = ".env"
				} else {
					fileName = fmt.Sprintf(".env.%s", envName)
				}

				file, err := os.Create(fileName)
				if err != nil {
					return err
				}

				keys := slices.Collect(maps.Keys(vars))
				slices.Sort(keys)

				for k, v := range vars {
					fmt.Fprintf(file, "%s=%v\n", k, v)
				}
				file.Close()
			}

			fmt.Println("Envs pulled")
			return nil
		},
	}
	return cmd
}

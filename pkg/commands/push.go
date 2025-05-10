package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"env0/pkg/client"
)

func pushCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "push [envName]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Push local environment files to the remote app",
		RunE: func(cmd *cobra.Command, args []string) error {
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
			json.Unmarshal(data, &cfg)
			fullAppName := fmt.Sprintf("%s/%s", cfg.OwnerName, cfg.AppName)

			// 3) Read .env files
			envs := make(map[string]map[string]interface{})
			var target *string
			if len(args) == 1 {
				target = &args[0]
				if *target == defaultTargetEnv {
					target = &defaultTargetEnvKey
				}
			}
			files, _ := os.ReadDir(".")
			for _, fi := range files {
				name := fi.Name()
				if name != ".env" && !strings.HasPrefix(name, ".env.") {
					continue
				}

				var envName string
				if name == ".env" {
					envName = ""
				} else {
					envName = strings.TrimPrefix(name, ".env.")
				}

				if target != nil && envName != *target {
					continue
				}

				content, _ := os.ReadFile(name)
				lines := strings.Split(string(content), "\n")
				vars := make(map[string]interface{})
				for _, line := range lines {
					if strings.TrimSpace(line) == "" {
						continue
					}
					parts := strings.SplitN(line, "=", 2)
					vars[parts[0]] = parts[1]
				}
				envs[envName] = vars
			}

			// if target set but not found, ensure empty
			if target != nil && envs[*target] == nil {
				envs[*target] = map[string]interface{}{}
			}

			// 4) Push via API
			c := client.New(token)
			if err := c.UpdateApp(context.Background(), fullAppName, envs); err != nil {
				fmt.Println(err.Error())
				return nil
			}
			fmt.Println("Envs pushed")
			return nil
		},
	}
	return cmd
}

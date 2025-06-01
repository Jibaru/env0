package scripts

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Jibaru/env0/pkg/client"
	"github.com/Jibaru/env0/pkg/logger"
)

type CloneInput struct {
	FullAppName string
}

// CloneFn represents a function that performs the clone operation
type CloneFn func(context.Context, CloneInput) error

// NewClone creates a new clone function with injected dependencies
func NewClone(c client.Client, logger logger.Logger) CloneFn {
	return func(ctx context.Context, input CloneInput) error {
		// 1) Check for existing config
		if _, err := os.Stat(".env0/config.json"); err == nil {
			return fmt.Errorf("there is an app already cloned")
		}

		// 2) Fetch envs from API
		envs, err := c.GetApp(ctx, input.FullAppName)
		if err != nil {
			return err
		}

		// 3) Write .env files
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

		// 4) Save local config
		parts := strings.SplitN(input.FullAppName, "/", 2)
		owner := parts[0]
		app := parts[1]

		if err := os.MkdirAll(".env0", 0755); err != nil {
			return err
		}

		cfgContent := fmt.Sprintf("{\"appName\": \"%s\", \"ownerName\": \"%s\"}", app, owner)
		if err := os.WriteFile(filepath.Join(".env0", "config.json"), []byte(cfgContent), 0644); err != nil {
			return err
		}

		logger.Printf("App cloned successfully")
		return nil
	}
}

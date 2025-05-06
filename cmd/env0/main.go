package main

import (
	"os"

	"github.com/spf13/cobra"

	"env0/pkg/commands"
)

func main() {
	rootCmd := &cobra.Command{Use: "env0"}
	commands.RegisterCommands(rootCmd)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

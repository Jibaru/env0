package commands

import (
	"github.com/spf13/cobra"
)

func RegisterCommands(root *cobra.Command) {
	root.AddCommand(
		signupCmd(),
		loginCmd(),
		initCmd(),
		cloneCmd(),
		pullCmd(),
		pushCmd(),
		addUserCmd(),
		delUserCmd(),
		versionCmd(),
	)
}

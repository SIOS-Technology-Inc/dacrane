package cmd

import (
	"dacrane/core"

	"github.com/spf13/cobra"
)

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "switch context",
	Long:  "switch context",
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		config := core.LoadContextConfig()
		config.Switch(name)
	},
}

func init() {
	contextCmd.AddCommand(switchCmd)
}

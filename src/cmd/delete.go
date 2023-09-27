package cmd

import (
	"dacrane/core"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete context",
	Long:  "delete context",
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		config := core.LoadContextConfig()
		config.Delete(name)
	},
}

func init() {
	contextCmd.AddCommand(deleteCmd)
}

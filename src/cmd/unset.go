package cmd

import (
	"dacrane/core"

	"github.com/spf13/cobra"
)

// unsetCmd represents the unset command
var unsetCmd = &cobra.Command{
	Use:   "unset",
	Short: "unset environment",
	Long:  "unset environment",
	Run: func(cmd *cobra.Command, args []string) {
		context := core.LoadContextConfig().CurrentContext()
		context.WriteEnv([]byte{})
	},
}

func init() {
	rootCmd.AddCommand(unsetCmd)
}

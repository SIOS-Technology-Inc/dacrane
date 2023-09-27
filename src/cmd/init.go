package cmd

import (
	"dacrane/core"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize dacrane project",
	Long:  "initialize dacrane project",
	Run: func(cmd *cobra.Command, args []string) {
		core.NewDefaultContextConfig().Init()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

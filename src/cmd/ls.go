package cmd

import (
	"dacrane/core"
	"fmt"

	"github.com/spf13/cobra"
)

// lsCmd represents the versions command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "show instance list",
	Long:  "show instance list",
	Run: func(cmd *cobra.Command, args []string) {
		config := core.LoadProjectConfig()
		list := config.PrettyList()
		fmt.Print(list)
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}

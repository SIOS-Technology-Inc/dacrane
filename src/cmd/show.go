package cmd

import (
	"dacrane/core"
	"fmt"

	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show context list",
	Long:  `show context list`,
	Run: func(cmd *cobra.Command, args []string) {
		config := core.LoadContextConfig()
		list := config.PrettyList()
		fmt.Print(list)
	},
}

func init() {
	contextCmd.AddCommand(showCmd)
}

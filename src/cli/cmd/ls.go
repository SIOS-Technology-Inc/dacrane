package cmd

import (
	"dacrane/cli/core/module"
	"dacrane/cli/core/repository"
	"fmt"

	"github.com/spf13/cobra"
)

// lsCmd represents the versions command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "show instance list",
	Long:  "show instance list",
	Run: func(cmd *cobra.Command, args []string) {
		instances := repository.LoadDocumentRepository()
		list := module.PrettyInstanceList(instances)
		fmt.Print(list)
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}

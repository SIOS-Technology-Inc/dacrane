package cmd

import (
	"fmt"

	"github.com/SIOS-Technology-Inc/dacrane/v0/src/core/module"
	"github.com/SIOS-Technology-Inc/dacrane/v0/src/core/repository"

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

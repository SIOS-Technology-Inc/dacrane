package cmd

import (
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for the contents indicated by the specific subcommand",
	Long:  "",
}

func init() {
	rootCmd.AddCommand(searchCmd)
}

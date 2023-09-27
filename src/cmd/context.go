package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// contextCmd represents the context command
var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "context management",
	Long:  "context management",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("context called")
	},
}

func init() {
	rootCmd.AddCommand(contextCmd)
}

package cmd

import (
	"dacrane/core"
	"os"

	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "set variables",
	Long:  "set variables",
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]
		context := core.LoadContextConfig().CurrentContext()
		data, err := os.ReadFile(file)
		if err != nil {
			panic(err)
		}
		context.WriteEnv(data)
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}

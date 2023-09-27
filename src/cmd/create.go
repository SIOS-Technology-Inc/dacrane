/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"dacrane/core"

	"github.com/spf13/cobra"
)

// contextCreateCmd represents the create command
var contextCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create context",
	Long:  "create context",
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		config := core.LoadContextConfig()
		config.Add(core.Context{
			Name: name,
		})
	},
}

func init() {
	contextCmd.AddCommand(contextCreateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

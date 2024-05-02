/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/SIOS-Technology-Inc/dacrane/v0/src/langserver"
	"github.com/spf13/cobra"
)

// serveCmd represents the dls command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start Language Server",
	Long:  `Start Language Server`,
	Run: func(cmd *cobra.Command, args []string) {
		langserver.Start()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dlsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dlsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

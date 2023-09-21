/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"dacrane/core/code"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		codeBytes, err := os.ReadFile("dacrane.yaml")
		if err != nil {
			panic(err)
		}

		code, err := code.ParseCode(codeBytes)
		if err != nil {
			panic(err)
		}

		envBytes, err := os.ReadFile(".env.yaml")
		if err != nil {
			panic(err)
		}
		data := map[string]any{}
		env := map[string]any{}
		yaml.Unmarshal(envBytes, &env)
		data["data"] = env

		for entity := range code {
			println(entity)
		}
	},
}

func init() {
	rootCmd.AddCommand(upCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// upCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

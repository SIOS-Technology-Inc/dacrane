/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"dacrane/core"
	"dacrane/core/code"
	"errors"
	"os"

	"github.com/spf13/cobra"
)

// versionsCmd represents the versions command
var versionsCmd = &cobra.Command{
	Use:   "versions",
	Short: "Search for the specific artifact versions",
	Long:  "",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires artifact name")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		codeBytes, err := os.ReadFile("dacrane.yaml")
		if err != nil {
			panic(err)
		}

		code, err := code.ParseCode(codeBytes)
		if err != nil {
			panic(err)
		}

		artifactCode := code.Find("artifact", name)

		artifactProvider := core.FindArtifactProvider(artifactCode["provider"].(string))

		err = artifactProvider.SearchVersions(artifactCode["parameters"].(map[string]any))

		if err != nil {
			panic(err)
		}
		println("search versions successfully!")
	},
}

func init() {
	searchCmd.AddCommand(versionsCmd)
}

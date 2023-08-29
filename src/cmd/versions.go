/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"dacrane/core"
	"dacrane/utils"
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

		codes, err := core.ParseCode(codeBytes)
		if err != nil {
			panic(err)
		}

		artifactCode := utils.Find(codes, func(code core.Code) bool {
			return code.Kind == "artifact" && code.Name == name
		})

		artifactProvider := core.FindArtifactProvider(artifactCode.Provider)

		err = artifactProvider.SearchVersions(artifactCode.Parameters)

		if err != nil {
			panic(err)
		}
		println("search versions successfully!")
	},
}

func init() {
	searchCmd.AddCommand(versionsCmd)
}

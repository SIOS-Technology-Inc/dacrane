package cmd

import (
	"dacrane/core"
	"errors"
	"os"

	"github.com/spf13/cobra"
)

// lsCmd represents the versions command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "show instance list",
	Long:  "show instance list",
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
	rootCmd.AddCommand(lsCmd)
}

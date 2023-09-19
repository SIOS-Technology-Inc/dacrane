package cmd

import (
	"dacrane/core"
	"dacrane/core/code"
	"errors"
	"os"

	"github.com/spf13/cobra"
)

// unpublishCmd represents the unpublish command
var unpublishCmd = &cobra.Command{
	Use:   "unpublish",
	Short: "Unpublish the specific artifact",
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

		err = artifactProvider.Unpublish(artifactCode["parameters"].(map[string]any))

		if err != nil {
			panic(err)
		}
		println("unpublish successfully!")
	},
}

func init() {
	rootCmd.AddCommand(unpublishCmd)
}

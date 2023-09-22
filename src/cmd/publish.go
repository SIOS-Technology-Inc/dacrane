package cmd

import (
	"dacrane/core"
	"dacrane/core/code"
	"errors"
	"os"

	"github.com/spf13/cobra"
)

// publishCmd represents the publish command
var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish the specific artifact",
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

		_, err = artifactProvider.Publish(artifactCode["parameters"].(map[string]any))

		if err != nil {
			panic(err)
		}
		println("publish successfully!")
	},
}

func init() {
	rootCmd.AddCommand(publishCmd)

}

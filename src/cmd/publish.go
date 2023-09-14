package cmd

import (
	"dacrane/core"
	"dacrane/core/code"
	"dacrane/utils"
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

		codes, err := code.ParseCode(codeBytes)
		if err != nil {
			panic(err)
		}

		artifactCode := utils.Find(codes, func(code code.RawCode) bool {
			return code.Kind == "artifact" && code.Name == name
		})

		artifactProvider := core.FindArtifactProvider(artifactCode.Provider)

		err = artifactProvider.Publish(artifactCode.Parameters)

		if err != nil {
			panic(err)
		}
		println("publish successfully!")
	},
}

func init() {
	rootCmd.AddCommand(publishCmd)

}

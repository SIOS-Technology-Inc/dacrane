package cmd

import (
	"dacrane/core"
	"dacrane/utils"
	"os"

	"github.com/spf13/cobra"
)

// unpublishCmd represents the unpublish command
var unpublishCmd = &cobra.Command{
	Use:   "unpublish",
	Short: "Unpublish the specific artifact",
	Long:  "",
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

		err = artifactProvider.Unpublish(artifactCode.Parameters)

		if err != nil {
			panic(err)
		}
		println("unpublish successfully!")
	},
}

func init() {
	rootCmd.AddCommand(unpublishCmd)
}

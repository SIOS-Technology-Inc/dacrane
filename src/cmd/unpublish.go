package cmd

import (
	"dacrane/core"
	"dacrane/utils"
	"fmt"
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
		dacranePath := fmt.Sprintf("%s/dacrane.yaml", workingDir)
		codeBytes, err := os.ReadFile(dacranePath)
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

		result, err := artifactProvider.Unpublish(workingDir, artifactCode.Parameters)

		if err != nil {
			println(string(result))
			panic(err)
		}
		println("unpublish successfully!")
	},
}

func init() {
	rootCmd.AddCommand(unpublishCmd)
}

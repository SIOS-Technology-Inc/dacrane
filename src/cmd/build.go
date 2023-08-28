package cmd

import (
	"dacrane/core"
	"dacrane/utils"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the specific artifact",
	Long:  "",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires artifact name")
		}
		return nil
	},
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

		buildCode := utils.Find(codes, func(code core.Code) bool {
			return code.Kind == "artifact" && code.Name == name
		})

		artifactProvider := core.FindArtifactProvider(buildCode.Provider)

		result, err := artifactProvider.Build(workingDir, buildCode.Parameters)

		if err != nil {
			println(string(result))
			panic(err)
		}
		println("build successfully!")
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}

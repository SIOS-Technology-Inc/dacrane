package cmd

import (
	"dacrane/core"
	"dacrane/core/code"
	"dacrane/utils"
	"errors"
	"os"

	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "A brief description of your command",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires resource name")
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

		resourceCode := utils.Find(codes, func(code code.RawCode) bool {
			return code.Kind == "resource" && code.Name == name
		})

		resourceProvider := core.FindResourceProvider(resourceCode.Provider)

		err = resourceProvider.Create(resourceCode.Parameters, resourceCode.Credentials)

		if err != nil {
			panic(err)
		}
		println("deploy successfully!")
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
}

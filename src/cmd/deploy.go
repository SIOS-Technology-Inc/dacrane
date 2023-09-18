package cmd

import (
	"dacrane/core"
	"dacrane/core/code"
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

		code, err := code.ParseCode(codeBytes)
		if err != nil {
			panic(err)
		}

		resourceCode := code.Find("resource", name)

		resourceProvider := core.FindResourceProvider(resourceCode["provider"].(string))

		err = resourceProvider.Create(resourceCode["parameters"].(map[string]any), resourceCode["credentials"].(map[string]any))

		if err != nil {
			panic(err)
		}
		println("deploy successfully!")
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
}

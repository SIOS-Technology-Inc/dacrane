package cmd

import (
	"dacrane/core"
	"dacrane/core/code"
	"errors"
	"os"

	"github.com/spf13/cobra"
)

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

		err = resourceProvider.Delete(resourceCode["parameters"].(map[string]any), resourceCode["credentials"].(map[string]any))

		if err != nil {
			panic(err)
		}
		println("destroy successfully!")
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)
}

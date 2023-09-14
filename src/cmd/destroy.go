package cmd

import (
	"dacrane/core"
	"dacrane/core/code"
	"dacrane/utils"
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

		codes, err := code.ParseCode(codeBytes)
		if err != nil {
			panic(err)
		}

		resourceCode := utils.Find(codes, func(code code.RawCode) bool {
			return code.Kind == "resource" && code.Name == name
		})

		resourceProvider := core.FindResourceProvider(resourceCode.Provider)

		err = resourceProvider.Delete(resourceCode.Parameters, resourceCode.Credentials)

		if err != nil {
			panic(err)
		}
		println("destroy successfully!")
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)
}

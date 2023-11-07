package cmd

import (
	"dacrane/core"
	"errors"
	"os"

	"github.com/spf13/cobra"
)

// destroyCmd represents the down command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "destroy resource and artifact",
	Long:  "destroy resource and artifact",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires instance name")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		instanceName := args[0]
		config := core.LoadProjectConfig()
		codeBytes, err := os.ReadFile("dacrane.yaml")
		if err != nil {
			panic(err)
		}
		modules := core.ParseModules(codeBytes)
		config.Destroy(instanceName, modules)
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)
}

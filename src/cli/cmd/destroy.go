package cmd

import (
	"dacrane/core"
	"dacrane/core/module"
	"dacrane/core/repository"
	"errors"

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
		instances := repository.LoadDocumentRepository()
		doc := instances.Find(instanceName)
		instance := module.NewInstanceFromDocument(doc)
		instance.Destroy(instanceName, &instances, core.Providers)
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)
}

package cmd

import (
	"errors"

	"github.com/SIOS-Technology-Inc/dacrane/v0/src/core/module"
	"github.com/SIOS-Technology-Inc/dacrane/v0/src/core/repository"

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
		for _, instanceName := range args {
			instances := repository.LoadDocumentRepository()
			doc := instances.Find(instanceName)
			instance := module.NewInstanceFromDocument(doc)
			instance.Destroy(instanceName, &instances)
		}
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)
}

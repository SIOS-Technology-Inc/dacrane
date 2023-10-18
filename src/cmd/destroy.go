package cmd

import (
	"dacrane/core"
	"dacrane/utils"
	"errors"
	"fmt"
	"strings"

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

		instance := config.GetInstance(instanceName)

		sortedModuleCalls := utils.Reverse(instance.Module.TopologicalSortedModuleCalls())
		for _, moduleCall := range sortedModuleCalls {
			state := instance.State["module"].(map[string]any)[moduleCall.Name]
			if state == nil {
				fmt.Printf("[%s (%s)] Skipped.\n", moduleCall.Name, moduleCall.Module)
				continue
			}

			modulePaths := strings.Split(moduleCall.Module, "/")
			kind := modulePaths[0]

			switch kind {
			case "resource":
				name := modulePaths[1]
				resourceProvider := core.FindResourceProvider(name)
				fmt.Printf("[%s (%s)] Deleting...\n", moduleCall.Name, moduleCall.Module)
				err := resourceProvider.Delete(state.(map[string]any))
				if err != nil {
					panic(err)
				}
				fmt.Printf("[%s (%s)] Deleted.\n", moduleCall.Name, moduleCall.Module)
			case "artifact":
				name := modulePaths[1]
				artifactProvider := core.FindArtifactProvider(name)
				fmt.Printf("[%s (%s)] Unpublish...\n", moduleCall.Name, moduleCall.Module)
				err := artifactProvider.Unpublish(state.(map[string]any))
				if err != nil {
					panic(err)
				}
				fmt.Printf("[%s (%s)] Unpublished.\n", moduleCall.Name, moduleCall.Module)
			case "data":

			}
			delete(instance.State["module"].(map[string]any), "")
			config.UpsertInstance(instance)
		}
		config.DeleteInstance(instanceName)
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)
}

package cmd

import (
	"dacrane/core"
	"dacrane/utils"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// applyCmd represents the up command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "create or update resource",
	Long:  "create or update resource",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("requires module name and instance name")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		targetModuleName := args[0]
		instanceName := args[1]
		config := core.LoadProjectConfig()
		codeBytes, err := os.ReadFile("dacrane.yaml")
		if err != nil {
			panic(err)
		}

		modules := core.ParseModules(codeBytes)
		module := utils.Find(modules, func(m core.Module) bool {
			return m.Name == targetModuleName
		})

		var argument map[string]any
		err = yaml.Unmarshal([]byte(argumentString), &argument)
		if err != nil {
			panic(err)
		}

		// for _, instanceName := range dependencies {
		// }

		state := map[string]any{
			"parameter": argument,
			"module":    map[string]any{},
		}

		instance := core.Instance{
			Name:   instanceName,
			Module: module,
			State:  state,
		}
		config.CreateInstance(instance)
		moduleCalls := module.TopologicalSortedModuleCalls()
		for _, moduleCall := range moduleCalls {
			fmt.Printf("[%s (%s)] Evaluating...\n", moduleCall.Name, moduleCall.Module)
			evaluatedModuleCall := moduleCall.Evaluate(state)
			fmt.Printf("[%s (%s)] Evaluated\n", moduleCall.Name, moduleCall.Module)

			modulePaths := strings.Split(evaluatedModuleCall.Module, "/")
			kind := modulePaths[0]

			switch kind {
			case "resource":
				name := modulePaths[1]
				resourceProvider := core.FindResourceProvider(name)
				fmt.Printf("[%s (%s)] Crating...\n", moduleCall.Name, moduleCall.Module)
				ret, err := resourceProvider.Create(evaluatedModuleCall.Argument.(map[string]any))
				if err != nil {
					panic(err)
				}
				fmt.Printf("[%s (%s)] Created.\n", moduleCall.Name, moduleCall.Module)
				state["module"].(map[string]any)[evaluatedModuleCall.Name] = ret
			case "artifact":
				name := modulePaths[1]
				artifactProvider := core.FindArtifactProvider(name)
				fmt.Printf("[%s (%s)] Building...\n", moduleCall.Name, moduleCall.Module)
				err = artifactProvider.Build(evaluatedModuleCall.Argument.(map[string]any))
				if err != nil {
					panic(err)
				}
				fmt.Printf("[%s (%s)] Built.\n", moduleCall.Name, moduleCall.Module)
				fmt.Printf("[%s (%s)] Publishing...\n", moduleCall.Name, moduleCall.Module)
				ret, err := artifactProvider.Publish(evaluatedModuleCall.Argument.(map[string]any))
				if err != nil {
					panic(err)
				}
				fmt.Printf("[%s (%s)] Published.\n", moduleCall.Name, moduleCall.Module)
				state["module"].(map[string]any)[evaluatedModuleCall.Name] = ret
			case "data":
				name := modulePaths[1]
				dataProvider := core.FindDataProvider(name)
				fmt.Printf("[%s (%s)] Reading...\n", moduleCall.Name, moduleCall.Module)
				ret, err := dataProvider.Get(evaluatedModuleCall.Argument.(map[string]any))
				if err != nil {
					panic(err)
				}
				fmt.Printf("[%s (%s)] Read.\n", moduleCall.Name, moduleCall.Module)
				state["module"].(map[string]any)[evaluatedModuleCall.Name] = ret
			}
			instance.State = state
			config.UpdateInstance(instance)
		}
	},
}

var argumentString = ""
var dependencies = map[string]string{}

func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.Flags().StringVarP(&argumentString, "argument", "a", "{}", "Argument")
	applyCmd.Flags().StringToStringVarP(&dependencies, "dependency", "d", map[string]string{}, "Argument")
}

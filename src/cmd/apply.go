package cmd

import (
	"dacrane/core"
	"dacrane/core/module"
	"dacrane/core/repository"
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
		moduleName := args[0]
		instanceName := args[1]
		instances := repository.LoadDocumentRepository()
		codeBytes, err := os.ReadFile("dacrane.yaml")
		if err != nil {
			panic(err)
		}

		modules := module.ParseModules(codeBytes)

		var argument map[string]any
		err = yaml.Unmarshal([]byte(argumentString), &argument)
		if err != nil {
			panic(err)
		}

		states := map[string]any{}
		for address, doc := range instances.Document() {
			instance := module.NewInstanceFromDocument(doc)
			if !strings.Contains(address, ".") {
				states[address] = instance.ToState(instances)
			}
		}

		evaluatedArg := module.Evaluate(argument, map[string]any{
			"instances": states,
		})

		moduleExists := utils.Contains(modules, func(module module.Module) bool {
			return module.Name == moduleName
		})
		if !moduleExists {
			panic(fmt.Sprintf("undefined module: %s", moduleName))
		}

		module := utils.Find(modules, func(module module.Module) bool {
			return module.Name == moduleName
		})

		module.Apply(instanceName, evaluatedArg, &instances, modules, core.Providers)
	},
}

var argumentString = ""

func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.Flags().StringVarP(&argumentString, "argument", "a", "{}", "Argument")
}

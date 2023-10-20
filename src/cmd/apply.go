package cmd

import (
	"dacrane/core"
	"errors"
	"os"

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
		config := core.LoadProjectConfig()
		codeBytes, err := os.ReadFile("dacrane.yaml")
		if err != nil {
			panic(err)
		}

		modules := core.ParseModules(codeBytes)

		var argument map[string]any
		err = yaml.Unmarshal([]byte(argumentString), &argument)
		if err != nil {
			panic(err)
		}

		config.Apply(instanceName, moduleName, argument, dependencies, modules)
	},
}

var argumentString = ""
var dependencies = map[string]string{}

func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.Flags().StringVarP(&argumentString, "argument", "a", "{}", "Argument")
	applyCmd.Flags().StringToStringVarP(&dependencies, "dependency", "d", map[string]string{}, "Argument")
}

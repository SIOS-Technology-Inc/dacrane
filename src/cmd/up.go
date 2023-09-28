package cmd

import (
	"dacrane/core"
	"dacrane/core/code"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "deploy resource and build artifact",
	Long:  "deploy resource and build artifact",
	Run: func(cmd *cobra.Command, args []string) {
		context := core.LoadContextConfig().CurrentContext()
		codeBytes, err := os.ReadFile("dacrane.yaml")
		if err != nil {
			panic(err)
		}

		code, err := code.ParseCode(codeBytes)
		if err != nil {
			panic(err)
		}

		envBytes := context.ReadEnv()
		env := map[string]any{}
		yaml.Unmarshal(envBytes, &env)

		inputs := map[string]any{}
		for k, v := range vars {
			fmt.Printf("%s = %s", k, v)
			inputs[k] = v
		}
		data := map[string]any{
			"config":   env,
			"arg":      inputs,
			"resource": map[string]any{},
			"artifact": map[string]any{},
			"data":     map[string]any{},
		}

		states := []map[string]any{}
		sortedEntities := code.TopologicalSort()
		for _, entity := range sortedEntities {
			fmt.Printf("[%s] Evaluating...\n", entity.Id())
			evaluatedEntity := entity.Evaluate(data)
			if evaluatedEntity == nil {
				fmt.Printf("[%s] Skipped.\n", entity.Id())
				continue
			}

			switch evaluatedEntity.Kind() {
			case "resource":
				resourceProvider := core.FindResourceProvider(entity.Provider())
				fmt.Printf("[%s] Crating...\n", entity.Id())
				ret, err := resourceProvider.Create(evaluatedEntity.Parameters())
				if err != nil {
					panic(err)
				}
				fmt.Printf("[%s] Created.\n", entity.Id())
				data["resource"].(map[string]any)[entity.Name()] = ret
			case "artifact":
				artifactProvider := core.FindArtifactProvider(evaluatedEntity.Provider())
				fmt.Printf("[%s] Building...\n", entity.Id())
				err = artifactProvider.Build(evaluatedEntity.Parameters())
				if err != nil {
					panic(err)
				}
				fmt.Printf("[%s] Built.\n", entity.Id())
				fmt.Printf("[%s] Publishing...\n", entity.Id())
				ret, err := artifactProvider.Publish(evaluatedEntity.Parameters())
				if err != nil {
					panic(err)
				}
				fmt.Printf("[%s] Published.\n", entity.Id())
				data["artifact"].(map[string]any)[entity.Name()] = ret
			case "data":
				dataProvider := core.FindDataProvider(evaluatedEntity.Provider())
				fmt.Printf("[%s] Reading...\n", entity.Id())
				ret, err := dataProvider.Get(entity.Parameters())
				if err != nil {
					panic(err)
				}
				fmt.Printf("[%s] Read.\n", entity.Id())
				data["data"].(map[string]any)[entity.Name()] = ret
			}
			states = append(states, evaluatedEntity)
			statesYaml := []byte{}
			for _, state := range states {
				stateYaml, e := yaml.Marshal(state)
				statesYaml = append(statesYaml, []byte("---\n")...)
				statesYaml = append(statesYaml, stateYaml...)
				if e != nil {
					panic(e)
				}
			}

			context.WriteState(statesYaml)
		}
	},
}

var vars = map[string]string{}

func init() {
	rootCmd.AddCommand(upCmd)
	upCmd.Flags().StringToStringVarP(&vars, "argument", "a", nil, "Argument")
}

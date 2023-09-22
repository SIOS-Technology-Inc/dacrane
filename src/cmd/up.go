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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		codeBytes, err := os.ReadFile("dacrane.yaml")
		if err != nil {
			panic(err)
		}

		code, err := code.ParseCode(codeBytes)
		if err != nil {
			panic(err)
		}

		envBytes, err := os.ReadFile(".env.yaml")
		if err != nil {
			panic(err)
		}
		env := map[string]any{}
		yaml.Unmarshal(envBytes, &env)
		data := map[string]any{
			"data":     env,
			"resource": map[string]any{},
			"artifact": map[string]any{},
		}

		sortedEntities := code.TopologicalSort()
		for _, entity := range sortedEntities {
			fmt.Printf("============ %s.%s ============\n", entity.Kind(), entity.Name())
			evaluatedEntity := entity.Evaluate(data)
			if evaluatedEntity == nil {
				println("No Resource Provided")
				continue
			}
			yaml, e := yaml.Marshal(evaluatedEntity)
			if e != nil {
				panic(e)
			}

			fmt.Println(string(yaml))
			switch evaluatedEntity.Kind() {
			case "resource":
				resourceProvider := core.FindResourceProvider(entity.Provider())
				ret, err := resourceProvider.Create(evaluatedEntity.Parameters())
				if err != nil {
					panic(err)
				}
				data["resource"].(map[string]any)[entity.Name()] = ret
			case "artifact":
				artifactProvider := core.FindArtifactProvider(evaluatedEntity.Provider())

				err = artifactProvider.Build(evaluatedEntity.Parameters())
				err = artifactProvider.Publish(evaluatedEntity.Parameters())
			case "data":

			}
		}
	},
}

func init() {
	rootCmd.AddCommand(upCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// upCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

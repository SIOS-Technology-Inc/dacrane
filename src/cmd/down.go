package cmd

import (
	"dacrane/core"
	"dacrane/core/code"
	"dacrane/utils"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "destroy resource and artifact",
	Long:  "destroy resource and artifact",
	Run: func(cmd *cobra.Command, args []string) {
		stateBytes, err := os.ReadFile(".dacrane.state.yaml")
		if err != nil {
			panic(err)
		}

		states, err := code.ParseCode(stateBytes)
		if err != nil {
			panic(err)
		}

		codeBytes, err := os.ReadFile("dacrane.yaml")
		if err != nil {
			panic(err)
		}

		dcode, err := code.ParseCode(codeBytes)
		if err != nil {
			panic(err)
		}

		sortedEntities := utils.Reverse(dcode.TopologicalSort())
		for _, entity := range sortedEntities {
			stateEntity := states.Find(entity.Kind(), entity.Name())
			if stateEntity == nil {
				println("No Resource Provided")
				continue
			}

			switch stateEntity.Kind() {
			case "resource":
				resourceProvider := core.FindResourceProvider(entity.Provider())
				fmt.Printf("[%s] Deleting...\n", entity.Id())
				err := resourceProvider.Delete(stateEntity.Parameters())
				if err != nil {
					panic(err)
				}
				fmt.Printf("[%s] Deleted.\n", entity.Id())
			case "artifact":
				artifactProvider := core.FindArtifactProvider(stateEntity.Provider())
				fmt.Printf("[%s] Unpublish...\n", entity.Id())
				err = artifactProvider.Unpublish(stateEntity.Parameters())
				if err != nil {
					panic(err)
				}
				fmt.Printf("[%s] Unpublished.\n", entity.Id())
			case "data":

			}
			states = utils.Filter(states, func(e code.Entity) bool {
				return e.Id() != entity.Id()
			})
			statesYaml := []byte{}
			for _, state := range states {
				stateYaml, e := yaml.Marshal(state)
				statesYaml = append(statesYaml, []byte("---\n")...)
				statesYaml = append(statesYaml, stateYaml...)
				if e != nil {
					panic(e)
				}
			}

			e := os.WriteFile(".dacrane.state.yaml", statesYaml, 0644)
			if e != nil {
				panic(e)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
}

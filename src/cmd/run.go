package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/SIOS-Technology-Inc/dacrane/v0/src/ast"
	"github.com/SIOS-Technology-Inc/dacrane/v0/src/parser"
	"github.com/spf13/cobra"
)

// runCmd represents the up command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run function",
	Long:  "run function",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires file name")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fileName := args[0]

		codeBytes, err := os.ReadFile(fileName)
		if err != nil {
			panic(err)
		}
		expr := parser.Parse(string(codeBytes))
		res, err := expr.Evaluate()
		var evalErr *ast.EvalError
		if errors.As(err, &evalErr) {
			os.Stderr.Write([]byte(fileName + ":" + evalErr.Position.String() + ":" + err.Error() + "\n"))
			return
		}
		if err != nil {
			os.Stderr.Write([]byte(err.Error() + "\n"))
			return
		}
		fmt.Println(res)
	},
}

var argumentString map[string]string

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringToStringVarP(&argumentString, "argument", "a", map[string]string{}, "Argument")
}

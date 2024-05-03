package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/SIOS-Technology-Inc/dacrane/v0/src/exception"
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
		tokens, err := parser.Lex(string(codeBytes))
		var codeErr *exception.CodeError
		if errors.As(err, &codeErr) {
			fmt.Fprintln(os.Stderr, codeErr.Pretty(fileName))
			return
		}
		expr := parser.Parse(tokens)
		res, err := expr.Evaluate()
		if errors.As(err, &codeErr) {
			fmt.Fprintln(os.Stderr, codeErr.Pretty(fileName))
			return
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
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

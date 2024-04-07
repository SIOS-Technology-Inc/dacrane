package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvaluate(t *testing.T) {
	var expr Expr = App{
		Func: "+",
		Args: []Expr{SInteger{Value: 1}, SInteger{Value: 2}},
	}
	result, err := expr.Evaluate(map[string]any{})
	assert.NoError(t, err)
	assert.Equal(t, 3, result)
}

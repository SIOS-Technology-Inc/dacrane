package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvaluate(t *testing.T) {
	var expr Expr = App{
		Func: "+",
		Args: []Expr{SInt{Value: 1}, SInt{Value: 2}},
	}
	result, err := expr.Evaluate([]Var{})
	assert.NoError(t, err)
	assert.Equal(t, 3, result)
}

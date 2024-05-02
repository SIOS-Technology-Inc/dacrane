package parser

import (
	"testing"

	"github.com/SIOS-Technology-Inc/dacrane/v0/src/ast"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	expr := Parse("1 + 1")
	assert.Equal(t, &ast.App{
		Func: "+",
		Args: []ast.Expr{
			&ast.SInt{Value: 1, Pos: ast.Position{Line: 0, Column: 0}},
			&ast.SInt{Value: 1, Pos: ast.Position{Line: 0, Column: 4}},
		},
		Pos: ast.Position{Line: 0, Column: 0},
	}, expr)
}

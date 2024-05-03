package parser

import (
	"testing"

	"github.com/SIOS-Technology-Inc/dacrane/v0/src/ast"
	"github.com/SIOS-Technology-Inc/dacrane/v0/src/locator"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tokens, err := Lex("1 + 1")
	assert.NoError(t, err)
	expr := Parse(tokens)
	assert.Equal(t, &ast.App{
		Func: "+",
		Args: []ast.Expr{
			&ast.SInt{
				Value: 1,
				Range: locator.Range{
					Start: locator.Position{Line: 0, Column: 0},
					End:   locator.Position{Line: 0, Column: 1},
				},
			},
			&ast.SInt{
				Value: 1,
				Range: locator.Range{
					Start: locator.Position{Line: 0, Column: 4},
					End:   locator.Position{Line: 0, Column: 5},
				},
			},
		},
		Range: locator.Range{
			Start: locator.Position{Line: 0, Column: 0},
			End:   locator.Position{Line: 0, Column: 5},
		},
	}, expr)
}

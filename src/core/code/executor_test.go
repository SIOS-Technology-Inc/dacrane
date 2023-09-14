package code

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	assert.Equal(t, NewExprParam(NewRefParam([]string{"a", "b", "c"})), ParseExpr("a.b.c"))
	assert.Equal(t, NewExprParam(NewRefParam([]string{"a"})), ParseExpr("a"))
}

func TestExec(t *testing.T) {
	assert.Equal(t, "OK", Exec([]string{"a", "b"}, map[string]any{
		"a": map[string]any{
			"b": "OK",
		},
	}))
	assert.Equal(t, "OK", Exec([]string{"a"}, map[string]any{
		"a": "OK",
	}))
}

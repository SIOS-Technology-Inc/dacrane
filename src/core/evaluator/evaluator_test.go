package evaluator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	assert.Equal(t, &Identifier{Name: "a.b.c"}, Parse("a.b.c"))
	assert.Equal(t, &Identifier{Name: "a"}, Parse("a"))
}

func TestEvaluate(t *testing.T) {
	expr := Parse("data.config.foo")
	v := Evaluate(expr, map[string]any{
		"data": map[string]any{
			"config": map[string]any{
				"foo": "OK",
			},
		},
	})
	assert.Equal(t, "OK", v)

	expr = Parse("1 + 2 + 4")
	v = Evaluate(expr, map[string]any{})
	assert.Equal(t, 7.0, v)

	expr = Parse(`data.config.foo == "OK"`)
	v = Evaluate(expr, map[string]any{
		"data": map[string]any{
			"config": map[string]any{
				"foo": "OK",
			},
		},
	})
	assert.Equal(t, true, v)

	expr = Parse(`{"low": "Basic", "high": "Premium" }`)
	v = Evaluate(expr, map[string]any{})
	assert.Equal(t, map[Expr]any{
		"low":  "Basic",
		"high": "Premium",
	}, v)

	expr = Parse(`{"low": "Basic", "high": "Premium" }["high"]`)
	v = Evaluate(expr, map[string]any{})
	assert.Equal(t, "Premium", v)
}

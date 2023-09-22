package code

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	assert.Equal(t, &Identifier{Name: "a.b.c"}, ParseExpr("a.b.c"))
	assert.Equal(t, &Identifier{Name: "a"}, ParseExpr("a"))
}

func TestDependency(t *testing.T) {
	codeStr := `
kind: resource
name: a
provider: foo
parameters:
  a: ${ resource.b }
  b: 1
---
kind: resource
name: b
parameters:
  a: 1
  b: 2
`
	code, e := ParseCode([]byte(codeStr))
	assert.NoError(t, e)
	entities := code.Dependency("resource.a")
	assert.Len(t, entities, 1)
	assert.Equal(t, "resource", entities[0]["kind"])
	assert.Equal(t, "b", entities[0]["name"])
}

func TestEvaluate(t *testing.T) {
	expr := ParseExpr("data.config.foo")
	v := EvaluateExprString(expr, map[string]any{
		"data": map[string]any{
			"config": map[string]any{
				"foo": "OK",
			},
		},
	})
	assert.Equal(t, "OK", v)

	expr = ParseExpr("1 + 2 + 4")
	v = EvaluateExprString(expr, map[string]any{})
	assert.Equal(t, 7.0, v)

	expr = ParseExpr(`data.config.foo == "OK"`)
	v = EvaluateExprString(expr, map[string]any{
		"data": map[string]any{
			"config": map[string]any{
				"foo": "OK",
			},
		},
	})
	assert.Equal(t, true, v)
}

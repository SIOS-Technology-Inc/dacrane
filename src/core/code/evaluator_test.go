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
	entity := code.Find("resource", "a")
	paths := references(entity)
	assert.Equal(t, []string{"resource.b"}, paths)
}

func TestEvaluate(t *testing.T) {
	expr := ParseExpr("data.config.foo")
	v := Evaluate(expr, map[string]Value{"data.config.foo": StringValue("OK")})
	assert.Equal(t, StringValue("OK"), v)

	expr = ParseExpr("1 + 2 + 4")
	v = Evaluate(expr, map[string]Value{})
	assert.Equal(t, NumberValue(7), v)

	expr = ParseExpr(`data.config.foo == "OK"`)
	v = Evaluate(expr, map[string]Value{"data.config.foo": StringValue("OK")})
	assert.Equal(t, BoolValue(true), v)
}

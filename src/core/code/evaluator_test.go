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
	assert.Equal(t, entity["kind"], "resource")
	assert.Equal(t, entity["name"], "a")
	paths := references(entity)
	assert.Equal(t, []string{"resource.b"}, paths)
}

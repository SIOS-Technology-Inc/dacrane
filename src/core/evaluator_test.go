package core

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
name: foo
parameter:
  type: object
  properties:
    a: { type: number }
    b: { type: string, default: latest }
modules:
- name: foo
  depends_on:
    - bar
  module: resource/baz
  argument:
    a: ${{ module.baz }}
    b: abc
- name: bar
  module: resource/qux
  argument:
    a: 123
    b: abc
`
	modules := ParseModules([]byte(codeStr))
	assert.Len(t, modules, 1)
	module := modules[0]
	assert.Equal(t, "foo", module.Name)
	assert.Equal(t, []string{"bar", "baz"}, module.FindModuleCall("foo").Dependency())
}

func TestTopologicalSortedModuleCalls(t *testing.T) {
	codeStr := `
name: abc
parameter:
  type: object
  properties:
    a: { type: number }
    b: { type: string, default: latest }
modules:
- name: foo
  module: resource/a
  depends_on:
    - bar
- name: bar
  module: resource/a
  argument:
    a: ${{ module.baz }}
- name: baz
  module: resource/a
`
	modules := ParseModules([]byte(codeStr))
	assert.Len(t, modules, 1)
	module := modules[0]
	assert.Equal(t, "abc", module.Name)
	assert.Equal(t, []ModuleCall{
		module.FindModuleCall("baz"),
		module.FindModuleCall("bar"),
		module.FindModuleCall("foo"),
	},
		module.TopologicalSortedModuleCalls())
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

	expr = ParseExpr(`{"low": "Basic", "high": "Premium" }`)
	v = EvaluateExprString(expr, map[string]any{})
	assert.Equal(t, map[Expr]any{
		"low":  "Basic",
		"high": "Premium",
	}, v)

	expr = ParseExpr(`{"low": "Basic", "high": "Premium" }["high"]`)
	v = EvaluateExprString(expr, map[string]any{})
	assert.Equal(t, "Premium", v)
}

package module

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDependency(t *testing.T) {
	codeStr := `
name: foo
modules:
- name: foo
  depends_on:
    - bar
  module: resource/baz
  arguments:
    a: ${{ baz }}
    b: abc
- name: bar
  module: resource/qux
  arguments:
    a: 123
    b: abc
- name: baz
  module: resource/qux
  arguments:
    a: 123
    b: abc
`
	modules := ParseModules([]byte(codeStr))
	assert.Len(t, modules, 1)
	module := modules[0]
	assert.Equal(t, "foo", module.Name)
	assert.Equal(t, []string{"bar", "baz"}, module.FindModuleCall("foo").Dependency(module.ModuleNames()))
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
    a: ${{ baz }}
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

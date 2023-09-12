package expression

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	assert.Equal(t, []string{"a", "b", "c"}, Parse("a.b.c"))
	assert.Equal(t, []string{"a"}, Parse("a"))
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

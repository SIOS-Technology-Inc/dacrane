package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFillDefault(t *testing.T) {
	filledDoc, err := FillDefault(
		map[string]any{
			"type":    "object",
			"default": map[string]any{},
			"properties": map[string]any{
				"a": map[string]any{
					"type":    "integer",
					"default": 123,
				},
				"b": map[string]any{
					"type":    "string",
					"default": "abc",
				},
				"c": map[string]any{
					"type": "string",
				},
				"d": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"a": map[string]any{
							"type":    "string",
							"default": "123",
						},
						"b": map[string]any{
							"type": "string",
						},
					},
				},
			},
		},
		map[string]any{
			"b": "def",
			"d": map[string]any{},
		},
	)

	assert.NoError(t, err)
	assert.Equal(t, map[string]any{
		"a": 123,
		"b": "def",
		"c": nil,
		"d": map[string]any{
			"a": "123",
			"b": nil,
		},
	}, filledDoc)
}

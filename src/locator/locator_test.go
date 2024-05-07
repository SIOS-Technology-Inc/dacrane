package locator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnion(t *testing.T) {
	r1 := Range{
		Start: Position{Line: 0, Column: 0},
		End:   Position{Line: 0, Column: 1},
	}
	r2 := Range{
		Start: Position{Line: 0, Column: 3},
		End:   Position{Line: 0, Column: 4},
	}
	assert.Equal(t, Range{
		Start: Position{Line: 0, Column: 0},
		End:   Position{Line: 0, Column: 4},
	}, r1.Union(r2))
}

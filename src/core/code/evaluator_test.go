package code

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	assert.Equal(t, &Identifier{Name: "a.b.c"}, ParseExpr("a.b.c"))
	assert.Equal(t, &Identifier{Name: "a"}, ParseExpr("a"))
}

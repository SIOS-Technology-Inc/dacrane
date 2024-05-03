package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexError(t *testing.T) {
	_, err := Lex("1 + \"")
	assert.Errorf(t, err, "UnknownTokenError")
}

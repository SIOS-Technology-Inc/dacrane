package exception

import (
	"fmt"

	"github.com/SIOS-Technology-Inc/dacrane/v0/src/locator"
	"github.com/macrat/simplexer"
)

type CodeError struct {
	Range   locator.Range
	Message string
}

func (e *CodeError) Error() string {
	return e.Message
}

func (e *CodeError) Pretty(file string) string {
	return fmt.Sprintf("%s:%s: %s", file, e.Range.String(), e.Message)
}

func NewCodeError(r locator.Range, m string) *CodeError {
	return &CodeError{
		Range:   r,
		Message: m,
	}
}

func NewCodeErrorFrom(t simplexer.Token, m string) *CodeError {
	return NewCodeError(
		locator.NewRangeFromToken(t),
		m,
	)
}

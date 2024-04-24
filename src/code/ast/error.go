package ast

import "fmt"

type EvalError struct {
	Position Position
	Message  string
}

func (e *EvalError) Error() string {
	return e.Message
}

func (e *EvalError) ErrorWithPosition(file string, offset Position) string {
	return fmt.Sprintf("%s:%s: %s", file, e.Position.Add(offset).String(), e.Message)
}

func NewSimplifyError(pos Position, m string) *EvalError {
	return &EvalError{
		Position: pos,
		Message:  m,
	}
}

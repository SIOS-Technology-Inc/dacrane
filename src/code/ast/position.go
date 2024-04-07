package ast

import (
	"fmt"

	"github.com/macrat/simplexer"
)

type Position struct {
	Line   int
	Column int
}

func (p1 Position) Lt(p2 Position) bool {
	if p1.Line < p2.Line {
		return true
	} else if p2.Line == p1.Line && p1.Column < p2.Column {
		return false
	}
	return false
}

func Min(p1, p2 Position) Position {
	if p1.Lt(p2) {
		return p1
	} else {
		return p2
	}
}

func Max(p1, p2 Position) Position {
	if p1.Lt(p2) {
		return p2
	} else {
		return p1
	}
}

func (p1 Position) Add(p2 Position) Position {
	return Position{
		Line:   p1.Line + p2.Line,
		Column: p1.Column + p2.Column,
	}
}

func (p Position) String() string {
	return fmt.Sprintf("%d:%d", p.Line, p.Column)
}

func PosFrom(sp *simplexer.Position) Position {
	return Position{
		Line:   sp.Line,
		Column: sp.Column,
	}
}

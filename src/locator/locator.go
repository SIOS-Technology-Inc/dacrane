package locator

import (
	"fmt"
	"strings"

	"github.com/macrat/simplexer"
)

type Position struct {
	Line   int
	Column int
}

type Range struct {
	Start Position
	End   Position
}

func (p1 Position) Lt(p2 Position) bool {
	if p1.Line < p2.Line {
		return true
	} else if p2.Line == p1.Line && p1.Column < p2.Column {
		return true
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

func (r1 Range) Union(r2 Range) Range {
	return Range{
		Start: Min(r1.Start, r2.Start),
		End:   Max(r1.End, r2.End),
	}
}

func (p Position) String() string {
	return fmt.Sprintf("%d:%d", p.Line+1, p.Column+1)
}
func (r Range) String() string {
	return fmt.Sprintf("%s-%s", r.Start.String(), r.End.String())
}

func positionFromSimpleLexerPosition(sp simplexer.Position) Position {
	return Position{
		Line:   sp.Line,
		Column: sp.Column,
	}
}

func EndPos(p Position, s string) Position {
	lines := strings.Split(s, "\n")
	lineShift := len(lines) - 1

	if lineShift == 0 {
		p.Column += len(lines[0])
	} else {
		p.Column = len(lines[len(lines)-1])
	}
	p.Line += lineShift

	return p
}

func NewRangeFromToken(t simplexer.Token) Range {
	start := positionFromSimpleLexerPosition(t.Position)
	end := EndPos(start, t.Literal)
	return Range{Start: start, End: end}
}

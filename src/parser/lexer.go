package parser

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/SIOS-Technology-Inc/dacrane/v0/src/ast"
	"github.com/macrat/simplexer"
)

type Lexer struct {
	lexer     *simplexer.Lexer
	lastToken *simplexer.Token
	result    ast.Expr
}

func NewLexer(reader io.Reader) *Lexer {
	l := simplexer.NewLexer(reader)

	l.TokenTypes = []simplexer.TokenType{
		simplexer.NewRegexpTokenType(INTEGER, `[0-9]+`),
		simplexer.NewRegexpTokenType(STRING, `"([^"]*)"`),
		simplexer.NewRegexpTokenType(ADD, `\+`),
	}

	return &Lexer{lexer: l}
}

func (l *Lexer) Lex(lval *yySymType) int {
	token, err := l.lexer.Scan()
	if err != nil {
		if e, ok := err.(simplexer.UnknownTokenError); ok {
			fmt.Fprintln(os.Stderr, e.Error()+":")
			fmt.Fprintln(os.Stderr, l.lexer.GetLastLine())
			fmt.Fprintln(os.Stderr, strings.Repeat(" ", e.Position.Column)+strings.Repeat("^", len(e.Literal)))
		} else {
			l.Error(err.Error())
		}
		os.Exit(1)
	}
	if token == nil {
		return -1
	}

	lval.token = token

	l.lastToken = token

	return int(token.Type.GetID())
}

func (l *Lexer) Error(e string) {
	fmt.Fprintln(os.Stderr, e+":")
	fmt.Fprintln(os.Stderr, l.lexer.GetLastLine())
	fmt.Fprintln(os.Stderr, strings.Repeat(" ", l.lastToken.Position.Column)+strings.Repeat("^", len(l.lastToken.Literal)))
}

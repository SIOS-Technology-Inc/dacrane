package code

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/macrat/simplexer"
)

type Lexer struct {
	lexer     *simplexer.Lexer
	lastToken *simplexer.Token
	result    ExprParam
}

func NewLexer(reader io.Reader) *Lexer {
	l := simplexer.NewLexer(reader)

	l.TokenTypes = []simplexer.TokenType{
		simplexer.NewRegexpTokenType(NUMBER, `[0-9]+\.[0-9]+`),
		simplexer.NewRegexpTokenType(STRING, `".*"`),
		simplexer.NewRegexpTokenType(BOOLEAN, `true|false`),
		simplexer.NewRegexpTokenType(NULL, `null`),
		simplexer.NewRegexpTokenType(IDENTIFIER, `[a-zA-Z0-9_-]+`),
		simplexer.NewRegexpTokenType(DOT, `\.`),
		simplexer.NewRegexpTokenType(COMMA, `,`),
		simplexer.NewRegexpTokenType(COLON, `:`),
		simplexer.NewRegexpTokenType(AND, `&&`),
		simplexer.NewRegexpTokenType(OR, `\|\|`),
		simplexer.NewRegexpTokenType(NOT, `!`),
		simplexer.NewRegexpTokenType(EQ, `==`),
		simplexer.NewRegexpTokenType(LT, `<`),
		simplexer.NewRegexpTokenType(LTE, `<=`),
		simplexer.NewRegexpTokenType(GT, `>`),
		simplexer.NewRegexpTokenType(GTE, `>=`),
		simplexer.NewRegexpTokenType(PRIORITY, `>>`),
		simplexer.NewRegexpTokenType(ADD, `\+`),
		simplexer.NewRegexpTokenType(SUB, `-`),
		simplexer.NewRegexpTokenType(MUL, `\*`),
		simplexer.NewRegexpTokenType(DIV, `/`),
		simplexer.NewRegexpTokenType(LBRACKET, `\(`),
		simplexer.NewRegexpTokenType(RBRACKET, `\)`),
		simplexer.NewRegexpTokenType(LSBRACKET, `\[`),
		simplexer.NewRegexpTokenType(RSBRACKET, `\]`),
		simplexer.NewRegexpTokenType(LCBRACKET, `\{`),
		simplexer.NewRegexpTokenType(RCBRACKET, `\}`),
		simplexer.NewRegexpTokenType(IF, `if`),
		simplexer.NewRegexpTokenType(THEN, `then`),
		simplexer.NewRegexpTokenType(ELSE, `else`),
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

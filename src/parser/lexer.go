package parser

import (
	"errors"
	"strings"

	"github.com/SIOS-Technology-Inc/dacrane/v0/src/ast"
	"github.com/SIOS-Technology-Inc/dacrane/v0/src/exception"
	"github.com/SIOS-Technology-Inc/dacrane/v0/src/locator"
	"github.com/macrat/simplexer"
)

type TokenIterationLexer struct {
	tokens    []*simplexer.Token
	lastIndex int
	lastToken *simplexer.Token
	error     error
	result    ast.Expr
}

func NewTokenIterationLexer(tokens []*simplexer.Token) *TokenIterationLexer {
	return &TokenIterationLexer{
		tokens:    tokens,
		lastIndex: -1,
		lastToken: nil,
		result:    nil,
	}
}

func (l *TokenIterationLexer) Lex(lval *yySymType) int {
	l.lastIndex++
	if l.lastIndex > len(l.tokens)-1 {
		return -1
	}
	token := l.tokens[l.lastIndex]
	lval.token = token
	l.lastToken = token

	return int(token.Type.GetID())
}

func (l *TokenIterationLexer) Error(e string) {
	l.error = exception.NewCodeError(locator.NewRangeFromToken(*l.lastToken), e)
}

func Lex(code string) ([]*simplexer.Token, error) {
	reader := strings.NewReader(code)
	l := simplexer.NewLexer(reader)

	l.TokenTypes = []simplexer.TokenType{
		simplexer.NewRegexpTokenType(INTEGER, `[0-9]+`),
		simplexer.NewRegexpTokenType(STRING, `"([^"]*)"`),
		simplexer.NewRegexpTokenType(ADD, `\+`),
	}

	tokens := []*simplexer.Token{}

	for {
		token, err := l.Scan()

		if err != nil {
			var unknownTokenError simplexer.UnknownTokenError
			if errors.As(err, &unknownTokenError) {
				return nil, exception.NewCodeError(
					locator.NewRangeFromToken(simplexer.Token{
						Literal:  unknownTokenError.Literal,
						Position: unknownTokenError.Position,
					}),
					err.Error(),
				)
			} else {
				return nil, err
			}
		}

		if token == nil {
			break
		}

		tokens = append(tokens, token)
	}

	return tokens, nil
}

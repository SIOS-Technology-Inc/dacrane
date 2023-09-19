package code

import (
	"bytes"
	"strings"

	"gopkg.in/yaml.v3"
)

func ParseExpr(exprStr string) Expr {
	lexer := NewLexer(strings.NewReader(exprStr))
	yyParse(lexer)
	return lexer.result
}

func ParseCode(codeBytes []byte) (Code, error) {
	r := bytes.NewReader(codeBytes)
	dec := yaml.NewDecoder(r)

	var code Code
	for {
		var entity map[string]any
		if dec.Decode(&entity) != nil {
			break
		}
		code = append(code, entity)
	}

	return code, nil
}

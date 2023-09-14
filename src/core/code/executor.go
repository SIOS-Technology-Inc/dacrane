package code

import (
	"bytes"
	"strings"

	"gopkg.in/yaml.v3"
)

func ParseExpr(pathStr string) ExprParam {
	lexer := NewLexer(strings.NewReader(pathStr))
	yyParse(lexer)
	return lexer.result
}

func Exec(path Path, objects map[string]any) string {
	var value any = objects
	for _, key := range path {
		value = value.(map[string]any)[key]
	}
	return value.(string)
}

func ParseCode(codeBytes []byte) ([]RawCode, error) {
	r := bytes.NewReader(codeBytes)
	dec := yaml.NewDecoder(r)

	var codes []RawCode
	for {
		var code RawCode
		if dec.Decode(&code) != nil {
			break
		}
		codes = append(codes, code)
	}

	return codes, nil
}

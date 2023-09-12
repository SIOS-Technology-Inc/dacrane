package expression

import "strings"

func Parse(pathStr string) Path {
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

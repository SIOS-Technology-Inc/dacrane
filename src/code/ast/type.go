package ast

import (
	"fmt"
	"strings"
)

type Type string

const (
	TNull   Type = "null"
	TInt    Type = "int"
	TFloat  Type = "float"
	TBool   Type = "bool"
	TString Type = "string"
	TSeq    Type = "sequence"
	TMap    Type = "map"
)

func Signature(name string, types ...Type) string {
	ts := []string{}
	for _, t := range types {
		ts = append(ts, string(t))
	}
	return fmt.Sprintf("%s(%s)", name, strings.Join(ts, ", "))
}

func ToType(a any) Type {
	switch a.(type) {
	case nil:
		return TNull
	case int:
		return TInt
	case float64:
		return TFloat
	case bool:
		return TBool
	case string:
		return TString
	case []any:
		return TSeq
	case map[any]any:
		return TMap
	default:
		panic("unexpected data type")
	}
}

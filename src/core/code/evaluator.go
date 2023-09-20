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

type Value interface{}

type NumberValue float64
type StringValue string
type BoolValue bool
type NullValue struct{}
type ListValue []Value
type MapValue map[string]Value

func Evaluate(expr Expr, env map[string]Value) Value {
	switch e := expr.(type) {
	case *Identifier:
		if val, ok := env[e.Name]; ok {
			return val
		}
		panic("Undefined variable: " + e.Name)
	case *BinaryExpr:
		left := Evaluate(e.Left, env)
		right := Evaluate(e.Right, env)
		switch e.Op.Type.GetID() {
		case ADD:
			return NumberValue(left.(NumberValue) + right.(NumberValue))
		case SUB:
			return NumberValue(left.(NumberValue) - right.(NumberValue))
		case MUL:
			return NumberValue(left.(NumberValue) * right.(NumberValue))
		case DIV:
			return NumberValue(left.(NumberValue) / right.(NumberValue))
		case EQ:
			println(left.(StringValue), right.(StringValue))
			return BoolValue(left == right)
		case LT:
			return BoolValue(left.(NumberValue) < right.(NumberValue))
		case LTE:
			return BoolValue(left.(NumberValue) <= right.(NumberValue))
		case GT:
			return BoolValue(left.(NumberValue) > right.(NumberValue))
		case GTE:
			return BoolValue(left.(NumberValue) >= right.(NumberValue))
		case AND:
			return BoolValue(left.(BoolValue) && right.(BoolValue))
		case OR:
			return BoolValue(left.(BoolValue) || right.(BoolValue))
		}
	case *UnaryExpr:
		val := Evaluate(e.Expr, env)
		switch e.Op.Type.GetID() {
		case SUB:
			return NumberValue(-val.(NumberValue))
		case NOT:
			return BoolValue(!val.(BoolValue))
		}
	case *IfExpr:
		condition := Evaluate(e.Condition, env)
		if condition.(BoolValue) {
			return Evaluate(e.Then, env)
		}
		return Evaluate(e.Else, env)
	case *List:
		var values []Value
		for _, item := range e.Items {
			values = append(values, Evaluate(item, env))
		}
		return values
	case *Map:
		kvMap := make(map[string]Value)
		for k, v := range e.KVs {
			kvMap[k] = Evaluate(v, env)
		}
		return kvMap
	case *App:
		panic("App node evaluation is not supported in this example.")
	case *Number:
		return NumberValue(e.Value)
	case *String:
		return StringValue(e.Value)
	case *Boolean:
		return BoolValue(e.Value)
	case *Null:
		return NullValue{}
	}

	panic("Unsupported expression type")
}

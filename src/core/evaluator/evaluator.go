package evaluator

import (
	"strings"
)

func Parse(exprStr string) Expr {
	lexer := NewLexer(strings.NewReader(exprStr))
	yyParse(lexer)
	return lexer.result
}

func Evaluate(expr Expr, data map[string]any) any {
	switch e := expr.(type) {
	case *Identifier:
		keys := strings.Split(e.Name, ".")
		var val any = data
		for _, key := range keys {
			if v, ok := val.(map[string]any)[key]; ok {
				val = v
			} else {
				val = nil
			}
		}
		return val
	case *Ref:
		m := Evaluate(e.Expr, data).(map[Expr]any)
		key := Evaluate(e.Key, data)
		return m[key]
	case *BinaryExpr:
		left := Evaluate(e.Left, data)
		right := Evaluate(e.Right, data)
		switch e.Op.Type.GetID() {
		case PRIORITY:
			if left != nil {
				return left
			} else {
				return right
			}
		case ADD:
			return left.(float64) + right.(float64)
		case SUB:
			return left.(float64) - right.(float64)
		case MUL:
			return left.(float64) * right.(float64)
		case DIV:
			return left.(float64) / right.(float64)
		case EQ:
			return left == right
		case NEQ:
			return left != right
		case LT:
			return left.(float64) < right.(float64)
		case LTE:
			return left.(float64) <= right.(float64)
		case GT:
			return left.(float64) > right.(float64)
		case GTE:
			return left.(float64) >= right.(float64)
		case AND:
			return left.(bool) && right.(bool)
		case OR:
			return left.(bool) || right.(bool)
		}
	case *UnaryExpr:
		val := Evaluate(e.Expr, data)
		switch e.Op.Type.GetID() {
		case SUB:
			return -val.(float64)
		case NOT:
			return !val.(bool)
		}
	case *IfExpr:
		condition := Evaluate(e.Condition, data)
		if condition.(bool) {
			return Evaluate(e.Then, data)
		}
		return Evaluate(e.Else, data)
	case *List:
		var values []any
		for _, item := range e.Items {
			values = append(values, Evaluate(item, data))
		}
		return values
	case *Map:
		kvMap := make(map[Expr]any)
		for k, v := range e.KVs {
			kvMap[Evaluate(k, data)] = Evaluate(v, data)
		}
		return kvMap
	case *App:
		panic("App node evaluation is not supported in this example.")
	case *Number:
		return e.Value
	case *String:
		return e.Value
	case *Boolean:
		return e.Value
	case *Null:
		return nil
	}
	panic("Unsupported expression type")
}

func GetReferences(expr Expr) []string {
	var names []string

	switch e := expr.(type) {
	case *Ref:
		if id, ok := e.Expr.(*Identifier); ok {
			keys := strings.Split(id.Name, ".")
			if keys[0] == "modules" {
				names = append(names, keys[1])
			}
		}
		names = append(names, GetReferences(e.Expr)...)
		names = append(names, GetReferences(e.Key)...)
	case *Identifier:
		keys := strings.Split(e.Name, ".")
		if keys[0] == "modules" {
			names = append(names, keys[1])
		}
	case *BinaryExpr:
		names = append(names, GetReferences(e.Left)...)
		names = append(names, GetReferences(e.Right)...)
	case *UnaryExpr:
		names = append(names, GetReferences(e.Expr)...)
	case *IfExpr:
		names = append(names, GetReferences(e.Condition)...)
		names = append(names, GetReferences(e.Then)...)
		names = append(names, GetReferences(e.Else)...)
	case *List:
		for _, item := range e.Items {
			names = append(names, GetReferences(item)...)
		}
	case *Map:
		for _, v := range e.KVs {
			names = append(names, GetReferences(v)...)
		}
	case *App:
		for _, param := range e.Params {
			names = append(names, GetReferences(param)...)
		}
	}

	return names
}

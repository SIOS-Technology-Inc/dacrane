package evaluator

import (
	"fmt"
	"regexp"
	"strings"
)

func Parse(exprStr string) Expr {
	lexer := NewLexer(strings.NewReader(exprStr))
	yyParse(lexer)
	return lexer.result
}

func Evaluate(expr Expr, data map[string]any) any {
	switch e := expr.(type) {
	case *Ref:
		m := Evaluate(e.Expr, data)
		switch m := m.(type) {
		case map[Expr]any:
			key := Evaluate(e.Key, data)
			return m[key]
		case map[string]any:
			key := Evaluate(e.Key, data).(string)
			return m[key]
		case []any:
			key := Evaluate(e.Key, data).(float64)
			return m[int(key)]
		case nil:
			key := Evaluate(e.Key, data).(string)
			return data[key]
		default:
			panic(fmt.Sprintf("Unsupported expression type: %T", m))
		}
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
		default:
			panic("Unsupported binary operator type")
		}
	case *UnaryExpr:
		val := Evaluate(e.Expr, data)
		switch e.Op.Type.GetID() {
		case SUB:
			return -val.(float64)
		case NOT:
			return !val.(bool)
		default:
			panic("Unsupported unary operator type")
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
	default:
		panic(fmt.Sprintf("Unsupported expression type: %T", e))
	}
}

func HasReferences(expr Expr, pattern string) bool {
	refs := CollectReferences(expr, pattern)
	return len(refs) > 0
}

func isStaticRef(ref *Ref) bool {
	switch ref.Key.(type) {
	case *String:
		switch ref.Expr.(type) {
		case *Ref:
			return isStaticRef(ref.Expr.(*Ref))
		default:
			return true
		}
	default:
		return false
	}
}

func collectRefKey(ref *Ref) string {
	switch ref.Key.(type) {
	case *String:
		switch ref.Expr.(type) {
		case *Ref:
			parent := collectRefKey(ref.Expr.(*Ref))
			key := ref.Key.(*String).Value
			return parent + "." + key
		default:
			return ref.Key.(*String).Value
		}
	default:
		panic("it is not static reference.")
	}
}

func CollectReferences(expr Expr, pattern string) []string {
	var names []string
	r, err := regexp.Compile(pattern)
	if err != nil {
		panic(err)
	}

	switch e := expr.(type) {
	case *Ref:
		if isStaticRef(e) {
			id := collectRefKey(e)
			if r.MatchString(id) {
				names = append(names, id)
			}
		} else {
			names = append(names, CollectReferences(e.Expr, pattern)...)
			names = append(names, CollectReferences(e.Key, pattern)...)
		}
	case *BinaryExpr:
		names = append(names, CollectReferences(e.Left, pattern)...)
		names = append(names, CollectReferences(e.Right, pattern)...)
	case *UnaryExpr:
		names = append(names, CollectReferences(e.Expr, pattern)...)
	case *IfExpr:
		names = append(names, CollectReferences(e.Condition, pattern)...)
		names = append(names, CollectReferences(e.Then, pattern)...)
		names = append(names, CollectReferences(e.Else, pattern)...)
	case *List:
		for _, item := range e.Items {
			names = append(names, CollectReferences(item, pattern)...)
		}
	case *Map:
		for _, v := range e.KVs {
			names = append(names, CollectReferences(v, pattern)...)
		}
	case *App:
		for _, param := range e.Params {
			names = append(names, CollectReferences(param, pattern)...)
		}
	}

	return names
}

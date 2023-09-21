package code

import (
	"bytes"
	"dacrane/utils"
	"reflect"
	"regexp"
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

func Evaluate(expr Expr, env map[string]any) any {
	switch e := expr.(type) {
	case *Identifier:
		keys := strings.Split(e.Name, ".")
		var val any = env
		for _, key := range keys {
			if v, ok := val.(map[string]any)[key]; ok {
				val = v
			} else {
				panic("Undefined variable: " + e.Name)
			}
		}
		return val
	case *BinaryExpr:
		left := Evaluate(e.Left, env)
		right := Evaluate(e.Right, env)
		switch e.Op.Type.GetID() {
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
		val := Evaluate(e.Expr, env)
		switch e.Op.Type.GetID() {
		case SUB:
			return -val.(float64)
		case NOT:
			return !val.(bool)
		}
	case *IfExpr:
		condition := Evaluate(e.Condition, env)
		if condition.(bool) {
			return Evaluate(e.Then, env)
		}
		return Evaluate(e.Else, env)
	case *List:
		var values []any
		for _, item := range e.Items {
			values = append(values, Evaluate(item, env))
		}
		return values
	case *Map:
		kvMap := make(map[string]any)
		for k, v := range e.KVs {
			kvMap[k] = Evaluate(v, env)
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

func (code Code) Find(kind string, name string) Entity {
	return utils.Find(code, func(e Entity) bool {
		return e["kind"].(string) == kind && e["name"].(string) == name
	})
}

func (code Code) Dependency(kind string, name string) []Entity {
	entity := code.Find(kind, name)
	paths := references(entity)
	var dependencies []Entity
	for _, path := range paths {
		identifiers := strings.Split(path, ".")
		kind := identifiers[0]
		name := identifiers[1]
		dependency := code.Find(kind, name)
		if dependency != nil {
			dependencies = append(dependencies, dependency)
		}
	}
	return dependencies
}

func references(raw map[string]any) []string {
	var paths []string
	for _, v := range raw {
		switch t := reflect.TypeOf(v); t.Kind() {
		case reflect.Map:
			paths = append(paths, references(v.(map[string]any))...)
		case reflect.String:
			r, e := regexp.Compile(`\$\{(.*?)\}`)
			if e != nil {
				panic(e)
			}
			res := r.FindAllStringSubmatch(v.(string), -1)
			for _, exprStr := range res {
				print(exprStr[1])
				expr := ParseExpr(exprStr[1])
				paths = append(paths, extractRefNames(expr)...)
			}
		default:
		}
	}
	return paths
}

func extractRefNames(expr Expr) []string {
	var names []string

	switch e := expr.(type) {
	case *Ref:
		if id, ok := e.Expr.(*Identifier); ok {
			names = append(names, id.Name)
		}
		names = append(names, extractRefNames(e.Expr)...)
		names = append(names, extractRefNames(e.Key)...)
	case *Identifier:
		names = append(names, e.Name)
	case *BinaryExpr:
		names = append(names, extractRefNames(e.Left)...)
		names = append(names, extractRefNames(e.Right)...)
	case *UnaryExpr:
		names = append(names, extractRefNames(e.Expr)...)
	case *IfExpr:
		names = append(names, extractRefNames(e.Condition)...)
		names = append(names, extractRefNames(e.Then)...)
		names = append(names, extractRefNames(e.Else)...)
	case *List:
		for _, item := range e.Items {
			names = append(names, extractRefNames(item)...)
		}
	case *Map:
		for _, v := range e.KVs {
			names = append(names, extractRefNames(v)...)
		}
	case *App:
		for _, param := range e.Params {
			names = append(names, extractRefNames(param)...)
		}
	}

	return names
}

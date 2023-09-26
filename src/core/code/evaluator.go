package code

import (
	"bytes"
	"dacrane/utils"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
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

func EvaluateExprString(expr Expr, data map[string]any) any {
	switch e := expr.(type) {
	case *Identifier:
		keys := strings.Split(e.Name, ".")
		var val any = data
		for _, key := range keys {
			if v, ok := val.(map[string]any)[key]; ok {
				val = v
			} else {
				panic("Undefined variable: " + e.Name)
			}
		}
		return val
	case *Ref:
		m := EvaluateExprString(e.Expr, data).(map[Expr]any)
		key := EvaluateExprString(e.Key, data)
		return m[key]
	case *BinaryExpr:
		left := EvaluateExprString(e.Left, data)
		right := EvaluateExprString(e.Right, data)
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
		val := EvaluateExprString(e.Expr, data)
		switch e.Op.Type.GetID() {
		case SUB:
			return -val.(float64)
		case NOT:
			return !val.(bool)
		}
	case *IfExpr:
		condition := EvaluateExprString(e.Condition, data)
		if condition.(bool) {
			return EvaluateExprString(e.Then, data)
		}
		return EvaluateExprString(e.Else, data)
	case *List:
		var values []any
		for _, item := range e.Items {
			values = append(values, EvaluateExprString(item, data))
		}
		return values
	case *Map:
		kvMap := make(map[Expr]any)
		for k, v := range e.KVs {
			kvMap[EvaluateExprString(k, data)] = EvaluateExprString(v, data)
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
		return e.Kind() == kind && e.Name() == name
	})
}

func (code Code) FindById(id string) Entity {
	return utils.Find(code, func(e Entity) bool {
		return e.Id() == id
	})
}

func (code Code) TopologicalSort() []Entity {
	g := simple.NewDirectedGraph()

	idToEntity := map[int64]Entity{}
	nodes := map[string]graph.Node{}
	for _, entity := range code {
		node := g.NewNode()
		nodes[entity.Id()] = node
		g.AddNode(node)
		idToEntity[node.ID()] = entity
	}

	for _, entity := range code {
		ds := code.Dependency(entity.Id())
		for _, d := range ds {
			g.SetEdge(g.NewEdge(nodes[d.Id()], nodes[entity.Id()]))
		}
	}

	sorted, err := topo.Sort(g)
	if err != nil {
		panic(err)
	}
	return utils.Map(sorted, func(node graph.Node) Entity {
		return idToEntity[node.ID()]
	})
}

func (code Code) Dependency(id string) []Entity {
	entity := code.FindById(id)
	paths := append(references(entity), entity.Dependencies()...)
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

func (entity Entity) Kind() string {
	return entity["kind"].(string)
}

func (entity Entity) Name() string {
	return entity["name"].(string)
}

func (entity Entity) Provider() string {
	return entity["provider"].(string)
}

func (entity Entity) Id() string {
	return fmt.Sprintf("%s.%s", entity.Kind(), entity.Name())
}

func (entity Entity) Parameters() map[string]any {
	return entity["parameters"].(map[string]any)
}

func (entity Entity) Dependencies() []string {
	if entity["dependencies"] == nil {
		return []string{}
	}
	paths := []string{}
	for _, d := range entity["dependencies"].([]any) {
		paths = append(paths, d.(string))
	}
	return paths
}

func (entity Entity) Evaluate(data map[string]any) Entity {
	m := EvaluateMap(entity.ToMap(), data)
	if m == nil {
		return nil
	}
	return Entity(m.(map[string]any))
}

func (entity Entity) ToMap() map[string]any {
	m := map[string]any{}
	for k, v := range entity {
		m[k] = v
	}
	return m
}

func EvaluateMap(prop any, data map[string]any) any {
	switch prop := prop.(type) {
	case string:
		single := isSingleExprString(prop)
		if single {
			r, e := regexp.Compile(`^\$\{\{(.*?)\}\}$`)
			if e != nil {
				panic(e)
			}
			exprStr := r.FindStringSubmatch(prop)[1]
			expr := ParseExpr(exprStr)
			return EvaluateExprString(expr, data)
		} else {
			return expandExpr(prop, data)
		}
	case map[string]any:
		prop, exists := evalIfProp(prop, data)
		if !exists {
			return nil
		}
		output := map[string]any{}
		for k, v := range prop {
			output[k] = EvaluateMap(v, data)
		}
		return output
	default:
		return prop
	}
}

func expandExpr(prop string, data map[string]any) string {
	r, e := regexp.Compile(`\$\{\{(.*?)\}\}`)
	if e != nil {
		panic(e)
	}
	return r.ReplaceAllStringFunc(prop, func(s string) string {
		exprStr := r.FindStringSubmatch(s)
		expr := ParseExpr(exprStr[1])
		v := EvaluateExprString(expr, data)
		return convertToString(v)
	})
}

func evalIfProp(prop map[string]any, data map[string]any) (map[string]any, bool) {
	if condition, ok := prop["if"]; ok {
		if !EvaluateMap(condition, data).(bool) {
			return nil, false
		}
	}
	delete(prop, "if")
	return prop, true
}

func isSingleExprString(s string) bool {
	r, e := regexp.Compile(`^\$\{\{.*?\}\}$`)
	if e != nil {
		panic(e)
	}
	return r.MatchString(s)
}

func convertToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", v)
	}
}

func references(raw map[string]any) []string {
	var paths []string
	for _, v := range raw {
		switch t := reflect.TypeOf(v); t.Kind() {
		case reflect.Map:
			paths = append(paths, references(v.(map[string]any))...)
		case reflect.String:
			r, e := regexp.Compile(`\$\{\{(.*?)\}\}`)
			if e != nil {
				panic(e)
			}
			res := r.FindAllStringSubmatch(v.(string), -1)
			for _, exprStr := range res {
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

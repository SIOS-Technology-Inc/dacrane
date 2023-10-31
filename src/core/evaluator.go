package core

import (
	"bytes"
	"dacrane/utils"
	"fmt"
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

func ParseModules(codeBytes []byte) []ModuleCode {
	r := bytes.NewReader(codeBytes)
	dec := yaml.NewDecoder(r)

	modules := []ModuleCode{}
	for {
		var module ModuleCode
		if err := dec.Decode(&module); err != nil {
			if err.Error() == "EOF" {
				break
			}
			panic(err)
		}
		modules = append(modules, module)
	}

	return modules
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
				val = nil
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

func (module ModuleCode) FindModuleCall(name string) ModuleCall {
	return utils.Find(module.ModuleCalls, func(mc ModuleCall) bool {
		return mc.Name == name
	})
}

func (module ModuleCode) TopologicalSortedModuleCalls() []ModuleCall {
	g := simple.NewDirectedGraph()

	idToName := map[int64]string{}
	nodes := map[string]graph.Node{}
	for _, moduleCall := range module.ModuleCalls {
		node := g.NewNode()
		nodes[moduleCall.Name] = node
		g.AddNode(node)
		idToName[node.ID()] = moduleCall.Name
	}

	for _, mc := range module.ModuleCalls {
		ds := mc.Dependency()
		for _, d := range ds {
			g.SetEdge(g.NewEdge(nodes[d], nodes[mc.Name]))
		}
	}

	sorted, err := topo.Sort(g)
	if err != nil {
		panic(err)
	}
	return utils.Map(sorted, func(node graph.Node) ModuleCall {
		return module.FindModuleCall(idToName[node.ID()])
	})
}

// returns dependency module name
func (mc ModuleCall) Dependency() []string {
	return append(mc.ExplicitDependency(), mc.ImplicitDependency()...)
}

func (mc ModuleCall) ExplicitDependency() []string {
	return mc.DependsOn
}

func (mc ModuleCall) ImplicitDependency() []string {
	var paths []string
	paths = append(paths, references(mc.Name)...)
	paths = append(paths, references(mc.Module)...)
	paths = append(paths, references(mc.Argument)...)
	return paths
}

func (mc ModuleCall) Evaluate(data map[string]any) *ModuleCall {

	mapMc := mc.toMap()

	evaluated := Evaluate(mapMc, data)

	if evaluated == nil {
		return nil
	}

	return toModuleCall(evaluated.(map[string]any))
}

func (mc ModuleCall) toMap() map[string]any {
	if mc.If == nil {
		mc.If = true
	}
	return map[string]any{
		"name":       mc.Name,
		"depends_on": mc.DependsOn,
		"module":     mc.Module,
		"argument":   mc.Argument,
		"if":         mc.If,
	}
}

func toModuleCall(mc map[string]any) *ModuleCall {
	var dependsOn []string
	if mc["depends_on"] == nil {
		dependsOn = []string{}
	} else {
		dependsOn = mc["depends_on"].([]string)
	}

	return &ModuleCall{
		Name:      mc["name"].(string),
		DependsOn: dependsOn,
		Module:    mc["module"].(string),
		Argument:  mc["argument"],
	}
}

func Evaluate(prop any, data map[string]any) any {
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
			output[k] = Evaluate(v, data)
		}
		return output
	case []any:
		output := []any{}
		for _, v := range prop {
			output = append(output, Evaluate(v, data))
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
		if !Evaluate(condition, data).(bool) {
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

func references(raw any) []string {
	switch raw := raw.(type) {
	case map[string]any:
		var paths []string
		for _, v := range raw {
			paths = append(paths, references(v)...)
		}
		return paths
	case []any:
		var paths []string
		for _, v := range raw {
			paths = append(paths, references(v)...)
		}
		return paths
	case string:
		r, e := regexp.Compile(`\$\{\{(.*?)\}\}`)
		if e != nil {
			panic(e)
		}
		res := r.FindAllStringSubmatch(raw, -1)
		var paths []string
		for _, exprStr := range res {
			expr := ParseExpr(exprStr[1])
			paths = append(paths, localModuleReferences(expr)...)
		}
		return paths
	default:
		return []string{}
	}
}

func localModuleReferences(expr Expr) []string {
	var names []string

	switch e := expr.(type) {
	case *Ref:
		if id, ok := e.Expr.(*Identifier); ok {
			keys := strings.Split(id.Name, ".")
			if keys[0] == "modules" {
				names = append(names, keys[1])
			}
		}
		names = append(names, localModuleReferences(e.Expr)...)
		names = append(names, localModuleReferences(e.Key)...)
	case *Identifier:
		keys := strings.Split(e.Name, ".")
		if keys[0] == "modules" {
			names = append(names, keys[1])
		}
	case *BinaryExpr:
		names = append(names, localModuleReferences(e.Left)...)
		names = append(names, localModuleReferences(e.Right)...)
	case *UnaryExpr:
		names = append(names, localModuleReferences(e.Expr)...)
	case *IfExpr:
		names = append(names, localModuleReferences(e.Condition)...)
		names = append(names, localModuleReferences(e.Then)...)
		names = append(names, localModuleReferences(e.Else)...)
	case *List:
		for _, item := range e.Items {
			names = append(names, localModuleReferences(item)...)
		}
	case *Map:
		for _, v := range e.KVs {
			names = append(names, localModuleReferences(v)...)
		}
	case *App:
		for _, param := range e.Params {
			names = append(names, localModuleReferences(param)...)
		}
	}

	return names
}

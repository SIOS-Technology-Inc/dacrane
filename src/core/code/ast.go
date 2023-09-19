package code

import (
	"dacrane/utils"
	"reflect"
	"regexp"
	"strings"

	"github.com/macrat/simplexer"
)

type Code []Entity

type Entity map[string]any

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
		dependencies = append(dependencies, dependency)
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

				expr := ParseExpr(exprStr[1])
				paths = append(paths, expr.(Identifier).Name)
			}
		default:
			panic("unexpected parameter type")
		}
	}
	return paths
}

// Expr represents an expression in the AST.
type Expr interface{}

// Number represents a numeric value.
type Number struct {
	Value float64
}

// String represents a string value.
type String struct {
	Value string
}

// Boolean represents a boolean value.
type Boolean struct {
	Value bool
}

// Null represents a null value.
type Null struct{}

// BinaryExpr represents a binary operation.
type BinaryExpr struct {
	Left  Expr
	Op    *simplexer.Token
	Right Expr
}

// UnaryExpr represents a unary operation.
type UnaryExpr struct {
	Op   *simplexer.Token
	Expr Expr
}

// IfExpr represents an if-then-else expression.
type IfExpr struct {
	Condition Expr
	Then      Expr
	Else      Expr
}

// List represents a list of expressions.
type List struct {
	Items []Expr
}

// Map represents a map of string to expression.
type Map struct {
	KVs map[string]Expr
}

// App represents a function application.
type App struct {
	Name   string
	Params []Expr
}

// Ref represents a reference to another expression.
type Ref struct {
	Expr Expr
	Key  Expr
}

// Identifier represents an identifier.
type Identifier struct {
	Name string
}

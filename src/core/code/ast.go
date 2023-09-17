package code

import (
	"reflect"
	"regexp"

	"github.com/macrat/simplexer"
)

type ParamType string

const (
	// Map            ParamType = "map"
	// Number         ParamType = "number"
	// String         ParamType = "string"
	Bool ParamType = "bool"
	// Null           ParamType = "null"
	StringWithExpr ParamType = "string_with_expr"
	// Expr           ParamType = "expr"
	// Ref            ParamType = "ref"
)

type RawCode struct {
	Kind        string         `yaml:"kind"`
	Name        string         `yaml:"name"`
	Provider    string         `yaml:"provider"`
	Parameters  map[string]any `yaml:"parameters"`
	Credentials map[string]any `yaml:"credentials"`
}

type Code struct {
	Kind        string
	Name        string
	Provider    string
	Parameters  MapParam
	Credentials MapParam
}

type Param interface {
	Type() ParamType
}

func NewParam(raw any) Param {
	switch t := reflect.TypeOf(raw); t.Kind() {
	case reflect.Map:
		return NewMapParam(raw.(map[string]any))
	case reflect.String:
		includesExpandExpr, e := regexp.MatchString("\\$\\{.*?\\}", raw.(string))
		if e != nil {
			panic(e)
		}
		if includesExpandExpr {
			return NewStringWithExprParam(raw.(string))
		} else {
			return NewStringParam(raw.(string))
		}
	default:
		panic("unexpected parameter type")
	}
}

// define map parameter
type MapParam struct {
	Param
	raw      map[string]any
	children map[string]Param
}

func NewMapParam(raw map[string]any) MapParam {
	return MapParam{
		raw: raw,
	}
}

func (MapParam) Type() ParamType {
	return ""
}

func (p MapParam) Get(key string) Param {
	return p.children[key]
}

// define number parameter
type StringParam struct {
	Param
	raw string
}

func NewStringParam(raw string) StringParam {
	return StringParam{
		raw: raw,
	}
}

func (StringParam) Type() ParamType {
	return ""
}

func (p StringParam) Get() string {
	return p.raw
}

// define expr parameter
type StringWithExprParam struct {
	Param
	children []Param
	exprRaws []string
	raw      string
}

func NewStringWithExprParam(raw string) StringWithExprParam {
	// TODO separate expression from string
	// r, e := regexp.Compile(`\$\{(.*?)\}`)
	// if e != nil {
	// 	panic(e)
	// }

	// res := r.FindAllStringSubmatch(raw, -1)

	return StringWithExprParam{
		children: []Param{},
		raw:      raw,
	}
}

func (StringWithExprParam) Type() ParamType {
	return StringWithExpr
}

func (p StringWithExprParam) Get(env map[string]string) string {
	// TODO resolve ref and env, generate string
	return ""
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

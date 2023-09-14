package code

import (
	"reflect"
	"regexp"
	"strings"
)

type ParamType string

const (
	Map            ParamType = "map"
	Number         ParamType = "number"
	String         ParamType = "string"
	Bool           ParamType = "bool"
	Null           ParamType = "null"
	StringWithExpr ParamType = "string_with_expr"
	Expr           ParamType = "expr"
	Ref            ParamType = "ref"
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
	return Map
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
	return String
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

// define expr parameter
type ExprParam struct {
	ref RefParam
}

func NewExprParam(ref RefParam) ExprParam {
	return ExprParam{
		ref: ref,
	}
}

func (ExprParam) Type() ParamType {
	return Expr
}

// define ref parameter
type RefParam struct {
	path Path
}

func NewRefParam(path Path) RefParam {
	return RefParam{
		path: path,
	}
}

func (RefParam) Type() ParamType {
	return Ref
}

// returns path string
func (p RefParam) Get() string {
	return strings.Join(p.path, ".")
}

type Identifier = string
type Path = []Identifier

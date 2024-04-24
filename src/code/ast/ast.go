package ast

import "fmt"

type Expr interface {
	Position() Position
	Simplify(vars map[string]any) (any, error)
	CollectVariables() []string
}

type Module struct {
	Name    string
	Imports []string
	Exports []string
	Assigns []Assign
}

type Assign struct {
	Name string
	Expr Expr
	Pos  Position
}

type Func struct {
	Params []Param
	Body   Expr
	Pos    Position
}

func (v Func) Position() (pos Position) {
	return v.Pos
}

func (v Func) Simplify(vars map[string]any) (any, error) {
	return func() {}, nil
}

func (v Func) CollectVariables() []string {
	return []string{}
}

type Param struct {
	Name string
	Pos  Position
}

// App represents a function application.
type App struct {
	Func string
	Args []Expr
	Pos  Position
}

func (v App) Position() (pos Position) {
	return v.Pos
}

func (v App) Simplify(vars map[string]any) (any, error) {
	ts := []Type{}
	vs := []any{}
	for _, arg := range v.Args {
		v, err := arg.Simplify(vars)
		if err != nil {
			return nil, err
		}
		vs = append(vs, v)
		ts = append(ts, ToType(v))
	}
	sign := Signature(v.Func, ts...)
	f, ok := FixtureFunctions[sign]
	if !ok {
		return nil, NewSimplifyError(v.Pos, fmt.Sprintf("%s is not defined", sign))
	}
	return f(vs)
}

func (v App) CollectVariables() []string {
	refs := []string{}
	for _, arg := range v.Args {
		refs = append(refs, arg.CollectVariables()...)
	}
	return refs
}

// Variable represents a reference to another expression.
type Variable struct {
	Name string
	Pos  Position
}

func (v Variable) Position() (pos Position) {
	return v.Pos
}

func (v Variable) Simplify(vars map[string]any) (any, error) {
	e, ok := vars[v.Name]
	if !ok {
		return nil, NewSimplifyError(v.Pos, fmt.Sprintf("%s is not defined", v.Name))
	}
	return e, nil
}

func (v Variable) CollectVariables() []string {
	return []string{v.Name}
}

// Ref represents a reference to another expression.
type Ref struct {
	Dict Expr
	Key  Expr
	Pos  Position
}

func (r Ref) Position() (pos Position) {
	return r.Pos
}

func (r Ref) Simplify(vars map[string]any) (any, error) {
	dict, err := r.Dict.Simplify(vars)
	if err != nil {
		return nil, err
	}
	key, err := r.Key.Simplify(vars)
	if err != nil {
		return nil, err
	}

	switch dict := dict.(type) {
	case []any:
		switch key := key.(type) {
		case int:
			return dict[key], nil
		default:
			return nil, NewSimplifyError(r.Pos, "the index of the sequence must be an integer")
		}
	case map[any]any:
		return dict[key], nil
	default:
		return nil, NewSimplifyError(r.Pos, "it is neither sequence nor mapping")
	}
}

func (r Ref) CollectVariables() []string {
	dictRefs := r.Dict.CollectVariables()
	keyRefs := r.Key.CollectVariables()
	return append(dictRefs, keyRefs...)
}

// Seq represents a seq of expressions.
type Seq struct {
	Value []Expr
	Pos   Position
}

func (s Seq) Position() (pos Position) {
	return s.Pos
}

func (s Seq) Simplify(vars map[string]any) (any, error) {
	SimplifySlice := []any{}
	for _, v := range s.Value {
		SimplifyValue, err := v.Simplify(vars)
		if err != nil {
			return nil, err
		}
		SimplifySlice = append(SimplifySlice, SimplifyValue)
	}
	return SimplifySlice, nil
}

func (s Seq) CollectVariables() []string {
	refs := []string{}
	for _, v := range s.Value {
		refs = append(refs, v.CollectVariables()...)
	}
	return refs
}

// Map represents a map of expression to expression.
type Map struct {
	Value map[Expr]Expr
	Pos   Position
}

func (m Map) Position() Position {
	return m.Pos
}

func (m Map) Simplify(vars map[string]any) (any, error) {
	SimplifyMap := map[any]any{}
	for k, v := range m.Value {
		SimplifyKey, err := k.Simplify(vars)
		if err != nil {
			return nil, err
		}
		SimplifyValue, err := v.Simplify(vars)
		if err != nil {
			return nil, err
		}
		SimplifyMap[SimplifyKey] = SimplifyValue
	}
	return SimplifyMap, nil
}

func (m Map) CollectVariables() []string {
	refs := []string{}
	for k, v := range m.Value {
		refs = append(refs, k.CollectVariables()...)
		refs = append(refs, v.CollectVariables()...)
	}
	return refs
}

// Scaler String represents string value
type SString struct {
	Value string
	Pos   Position
}

func (v SString) Position() Position {
	return v.Pos
}

func (v SString) Simplify(map[string]any) (any, error) {
	return v.Value, nil
}

func (SString) CollectVariables() []string {
	return []string{}
}

// Scaler Float represents floating point value
type SInteger struct {
	Value int
	Pos   Position
}

func (v SInteger) Position() Position {
	return v.Pos
}

func (v SInteger) Simplify(map[string]any) (any, error) {
	return v.Value, nil
}

func (SInteger) CollectVariables() []string {
	return []string{}
}

// Scaler Float represents floating point value
type SFloat struct {
	Value float64
	Pos   Position
}

func (v SFloat) Position() Position {
	return v.Pos
}

func (v SFloat) Simplify(map[string]any) (any, error) {
	return v.Value, nil
}

func (SFloat) CollectVariables() []string {
	return []string{}
}

// Scaler Boolean represents bool value
type SBool struct {
	Value bool
	Pos   Position
}

func (v SBool) Position() Position {
	return v.Pos
}

func (v SBool) Simplify(map[string]any) (any, error) {
	return v.Value, nil
}

func (SBool) CollectVariables() []string {
	return []string{}
}

// Scaler Null represents null value
type SNull struct {
	Pos Position
}

func (v SNull) Position() Position {
	return v.Pos
}

func (SNull) Simplify(map[string]any) (any, error) {
	return nil, nil
}

func (SNull) CollectVariables() []string {
	return []string{}
}

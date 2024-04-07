package ast

import "fmt"

type Expr interface {
	Position() Position
	Evaluate(vars map[string]Expr) (any, error)
	CollectVariables() []string
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

func (v App) Evaluate(vars map[string]Expr) (any, error) {
	ts := []Type{}
	vs := []any{}
	for _, arg := range v.Args {
		v, err := arg.Evaluate(vars)
		if err != nil {
			return nil, err
		}
		vs = append(vs, v)
		ts = append(ts, ToType(v))
	}
	sign := Signature(v.Func, ts...)
	f, ok := FixtureFunctions[sign]
	if !ok {
		return nil, NewEvaluateError(v.Pos, fmt.Sprintf("%s is not defined", sign))
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

func (v Variable) Evaluate(vars map[string]Expr) (any, error) {
	e, ok := vars[v.Name]
	if !ok {
		return nil, NewEvaluateError(v.Pos, fmt.Sprintf("%s is not defined", v.Name))
	}
	return e.Evaluate(vars)
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

func (r Ref) Evaluate(vars map[string]Expr) (any, error) {
	dict, err := r.Dict.Evaluate(vars)
	if err != nil {
		return nil, err
	}
	key, err := r.Key.Evaluate(vars)
	if err != nil {
		return nil, err
	}

	switch dict := dict.(type) {
	case []any:
		switch key := key.(type) {
		case int:
			return dict[key], nil
		default:
			return nil, NewEvaluateError(r.Pos, "the index of the sequence must be an integer")
		}
	case map[any]any:
		return dict[key], nil
	default:
		return nil, NewEvaluateError(r.Pos, "it is neither sequence nor mapping")
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

func (s Seq) Evaluate(vars map[string]Expr) (any, error) {
	EvaluateSlice := []any{}
	for _, v := range s.Value {
		EvaluateValue, err := v.Evaluate(vars)
		if err != nil {
			return nil, err
		}
		EvaluateSlice = append(EvaluateSlice, EvaluateValue)
	}
	return EvaluateSlice, nil
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

func (m Map) Evaluate(vars map[string]Expr) (any, error) {
	EvaluateMap := map[any]any{}
	for k, v := range m.Value {
		EvaluateKey, err := k.Evaluate(vars)
		if err != nil {
			return nil, err
		}
		EvaluateValue, err := v.Evaluate(vars)
		if err != nil {
			return nil, err
		}
		EvaluateMap[EvaluateKey] = EvaluateValue
	}
	return EvaluateMap, nil
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

func (v SString) Evaluate(map[string]Expr) (any, error) {
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

func (v SInteger) Evaluate(map[string]Expr) (any, error) {
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

func (v SFloat) Evaluate(map[string]Expr) (any, error) {
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

func (v SBool) Evaluate(map[string]Expr) (any, error) {
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

func (SNull) Evaluate(map[string]Expr) (any, error) {
	return nil, nil
}

func (SNull) CollectVariables() []string {
	return []string{}
}

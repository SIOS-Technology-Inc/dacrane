package ast

import "fmt"

type Expr interface {
	Position() Position
	Eval(vars map[string]Expr) (any, error)
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

func (v App) Eval(vars map[string]Expr) (any, error) {
	ts := []Type{}
	vs := []any{}
	for _, arg := range v.Args {
		v, err := arg.Eval(vars)
		if err != nil {
			return nil, err
		}
		vs = append(vs, v)
		ts = append(ts, ToType(v))
	}
	sign := Signature(v.Func, ts...)
	f, ok := FixtureFunctions[sign]
	if !ok {
		return nil, NewEvalError(v.Pos, fmt.Sprintf("%s is not defined", sign))
	}
	return f(vs)
}

// Variable represents a reference to another expression.
type Variable struct {
	Name string
	Pos  Position
}

func (v Variable) Position() (pos Position) {
	return v.Pos
}

func (v Variable) Eval(vars map[string]Expr) (any, error) {
	e, ok := vars[v.Name]
	if !ok {
		return nil, NewEvalError(v.Pos, fmt.Sprintf("%s is not defined", v.Name))
	}
	return e.Eval(vars)
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

func (r Ref) Eval(vars map[string]Expr) (any, error) {
	dict, err := r.Dict.Eval(vars)
	if err != nil {
		return nil, err
	}
	key, err := r.Key.Eval(vars)
	if err != nil {
		return nil, err
	}

	switch dict := dict.(type) {
	case []any:
		switch key := key.(type) {
		case int:
			return dict[key], nil
		default:
			return nil, NewEvalError(r.Pos, "the index of the sequence must be an integer")
		}
	case map[any]any:
		return dict[key], nil
	default:
		return nil, NewEvalError(r.Pos, "it is neither sequence nor mapping")
	}
}

// Seq represents a seq of expressions.
type Seq struct {
	Value []Expr
	Pos   Position
}

func (s Seq) Position() (pos Position) {
	return s.Pos
}

func (s Seq) Eval(vars map[string]Expr) (any, error) {
	evalSlice := []any{}
	for _, v := range s.Value {
		evalValue, err := v.Eval(vars)
		if err != nil {
			return nil, err
		}
		evalSlice = append(evalSlice, evalValue)
	}
	return evalSlice, nil
}

// Map represents a map of expression to expression.
type Map struct {
	Value map[Expr]Expr
	Pos   Position
}

func (m Map) Position() Position {
	return m.Pos
}

func (m Map) Eval(vars map[string]Expr) (any, error) {
	evalMap := map[any]any{}
	for k, v := range m.Value {
		evalKey, err := k.Eval(vars)
		if err != nil {
			return nil, err
		}
		evalValue, err := v.Eval(vars)
		if err != nil {
			return nil, err
		}
		evalMap[evalKey] = evalValue
	}
	return evalMap, nil
}

// Scaler String represents string value
type SString struct {
	Value string
	Pos   Position
}

func (v SString) Position() Position {
	return v.Pos
}

func (v SString) Eval(map[string]Expr) (any, error) {
	return v.Value, nil
}

// Scaler Float represents floating point value
type SInteger struct {
	Value int
	Pos   Position
}

func (v SInteger) Position() Position {
	return v.Pos
}

func (v SInteger) Eval(map[string]Expr) (any, error) {
	return v.Value, nil
}

// Scaler Float represents floating point value
type SFloat struct {
	Value float64
	Pos   Position
}

func (v SFloat) Position() Position {
	return v.Pos
}

func (v SFloat) Eval(map[string]Expr) (any, error) {
	return v.Value, nil
}

// Scaler Boolean represents bool value
type SBool struct {
	Value bool
	Pos   Position
}

func (v SBool) Position() Position {
	return v.Pos
}

func (v SBool) Eval(map[string]Expr) (any, error) {
	return v.Value, nil
}

// Scaler Null represents null value
type SNull struct {
	Pos Position
}

func (v SNull) Position() Position {
	return v.Pos
}

func (SNull) Eval(map[string]Expr) (any, error) {
	return nil, nil
}

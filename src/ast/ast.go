package ast

import "fmt"

type Expr interface {
	Position() Position
	Evaluate() (any, error)
	Infer() (Type, error)
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

func (v App) Evaluate() (any, error) {
	ts := []Type{}
	args := []any{}
	for _, arg := range v.Args {
		t, err := arg.Infer()
		if err != nil {
			return nil, err
		}
		ts = append(ts, t)

		arg, err := arg.Evaluate()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}
	f := FindFixtureFunctions(v.Func, ts)
	if f == nil {
		return nil, NewSimplifyError(v.Pos, fmt.Sprintf("%s(%s) is not defined", v.Func, ArgsType(ts).String()))
	}
	return f.Function(args)
}

func (v App) Infer() (Type, error) {
	ts := []Type{}
	for _, arg := range v.Args {
		t, err := arg.Infer()
		if err != nil {
			return nil, err
		}
		ts = append(ts, t)
	}
	f := FindFixtureFunctions(v.Func, ts)
	if f == nil {
		return nil, NewSimplifyError(v.Pos, fmt.Sprintf("%s is not defined", ArgsType(ts).String()))
	}
	return f.Type.Returns, nil
}

// Scaler Integer represents integer value
type SInt struct {
	Value int
	Pos   Position
}

func (v SInt) Position() Position {
	return v.Pos
}

func (v SInt) Evaluate() (any, error) {
	return v.Value, nil
}

func (v SInt) Infer() (Type, error) {
	return TInt, nil
}

// Scaler String represents string value
type SString struct {
	Value string
	Pos   Position
}

func (v SString) Position() Position {
	return v.Pos
}

func (v SString) Evaluate() (any, error) {
	return v.Value, nil
}

func (v SString) Infer() (Type, error) {
	return TString, nil
}

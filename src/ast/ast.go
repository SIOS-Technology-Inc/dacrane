package ast

import (
	"fmt"

	"github.com/SIOS-Technology-Inc/dacrane/v0/src/exception"
	"github.com/SIOS-Technology-Inc/dacrane/v0/src/locator"
)

type Module struct {
	Vars []Var
}

type Expr interface {
	GetRange() locator.Range
	Evaluate(vars []Var) (any, error)
	Infer(vars []Var) (Type, error)
}

func (m Module) FindVar(name string) (Var, bool) {
	return findVar(name, m.Vars)
}

func findVar(name string, Vars []Var) (Var, bool) {
	for _, v := range Vars {
		if v.Name == name {
			return v, true
		}
	}
	return Var{}, false
}

// Var represents a variable assignment.
type Var struct {
	Name  string
	Expr  Expr
	Range locator.Range
}

// Ref represents a reference of variable.
type Ref struct {
	Name  string
	Range locator.Range
}

func (r Ref) GetRange() locator.Range {
	return r.Range
}

func (r Ref) Evaluate(vars []Var) (any, error) {
	v, ok := findVar(r.Name, vars)
	if !ok {
		return nil, exception.NewCodeError(r.Range, fmt.Sprintf("%s is not defined", r.Name))
	}
	return v.Expr.Evaluate(vars)
}

func (r Ref) Infer(vars []Var) (Type, error) {
	v, ok := findVar(r.Name, vars)
	if !ok {
		return nil, exception.NewCodeError(r.Range, fmt.Sprintf("%s is not defined", r.Name))
	}
	return v.Expr.Infer(vars)
}

// App represents a function application.
type App struct {
	Func  string
	Args  []Expr
	Range locator.Range
}

func (v App) GetRange() locator.Range {
	return v.Range
}

func (v App) Evaluate(vars []Var) (any, error) {
	ts := []Type{}
	args := []any{}
	for _, arg := range v.Args {
		t, err := arg.Infer(vars)
		if err != nil {
			return nil, err
		}
		ts = append(ts, t)

		arg, err := arg.Evaluate(vars)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}
	f := FindFixtureFunctions(v.Func, ts)
	if f == nil {
		return nil, exception.NewCodeError(v.Range, fmt.Sprintf("%s(%s) is not defined", v.Func, ArgsType(ts).String()))
	}
	return f.Function(args)
}

func (v App) Infer(vars []Var) (Type, error) {
	ts := []Type{}
	for _, arg := range v.Args {
		t, err := arg.Infer(vars)
		if err != nil {
			return nil, err
		}
		ts = append(ts, t)
	}
	f := FindFixtureFunctions(v.Func, ts)
	if f == nil {
		return nil, exception.NewCodeError(v.Range, fmt.Sprintf("%s is not defined", ArgsType(ts).String()))
	}
	return f.Type.Returns, nil
}

// Scaler Integer represents integer value
type SInt struct {
	Value int
	Range locator.Range
}

func (v SInt) GetRange() locator.Range {
	return v.Range
}

func (v SInt) Evaluate([]Var) (any, error) {
	return v.Value, nil
}

func (v SInt) Infer([]Var) (Type, error) {
	return TInt, nil
}

// Scaler String represents string value
type SString struct {
	Value string
	Range locator.Range
}

func (v SString) GetRange() locator.Range {
	return v.Range
}

func (v SString) Evaluate([]Var) (any, error) {
	return v.Value, nil
}

func (v SString) Infer([]Var) (Type, error) {
	return TString, nil
}

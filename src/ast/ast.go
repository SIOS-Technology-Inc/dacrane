package ast

import (
	"fmt"

	"github.com/SIOS-Technology-Inc/dacrane/v0/src/exception"
	"github.com/SIOS-Technology-Inc/dacrane/v0/src/locator"
)

type Expr interface {
	GetRange() locator.Range
	Evaluate() (any, error)
	Infer() (Type, error)
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
		return nil, exception.NewCodeError(v.Range, fmt.Sprintf("%s(%s) is not defined", v.Func, ArgsType(ts).String()))
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

func (v SInt) Evaluate() (any, error) {
	return v.Value, nil
}

func (v SInt) Infer() (Type, error) {
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

func (v SString) Evaluate() (any, error) {
	return v.Value, nil
}

func (v SString) Infer() (Type, error) {
	return TString, nil
}

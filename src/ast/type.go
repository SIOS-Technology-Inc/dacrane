package ast

import (
	"fmt"
	"strings"
)

type Type interface {
	String() string
}

type DataType string
type ArgsType []Type
type FuncType struct {
	Arguments ArgsType
	Returns   Type
}

const (
	TInt    DataType = "int"
	TString DataType = "string"
)

func (ft FuncType) String() string {
	return fmt.Sprintf("(%s): %s", ft.Arguments.String(), ft.Returns.String())
}

func (t DataType) String() string {
	return string(t)
}

func (args ArgsType) String() string {
	ts := []string{}
	for _, t := range args {
		ts = append(ts, t.String())
	}
	return strings.Join(ts, ", ")
}

func Equal(t1, t2 Type) bool {
	return t1.String() == t2.String()
}

func (ft FuncType) Applicable(argTypes ArgsType) (bool, error) {
	if len(ft.Arguments) != len(argTypes) {
		return false, fmt.Errorf("wrong number of arguments")
	}
	for i, t := range ft.Arguments {
		if !Equal(t, argTypes[i]) {
			return false, fmt.Errorf("wrong type argument")
		}
	}
	return true, nil
}

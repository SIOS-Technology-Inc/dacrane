package ast

type fixtureFunction struct {
	Name     string
	Type     FuncType
	Function func([]any) (any, error)
}

var fixtureFunctions = []fixtureFunction{
	{
		Name: "+",
		Type: FuncType{
			Arguments: []Type{TInt, TInt},
			Returns:   TInt,
		},
		Function: func(a []any) (any, error) { return a[0].(int) + a[1].(int), nil },
	},
	{
		Name: "+",
		Type: FuncType{
			Arguments: []Type{TString, TString},
			Returns:   TString,
		},
		Function: func(a []any) (any, error) { return a[0].(string) + a[1].(string), nil },
	},
}

func FindFixtureFunctions(name string, args ArgsType) *fixtureFunction {
	for _, v := range fixtureFunctions {
		ok, _ := v.Type.Applicable(args)
		if v.Name == name && ok {
			return &v
		}
	}
	return nil
}

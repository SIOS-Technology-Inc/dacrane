package ast

import (
	"encoding/base64"
)

var FixtureFunctions = map[string]func([]any) (any, error){
	Signature("==", TInt, TInt):        func(a []any) (any, error) { return a[0].(int) == a[1].(int), nil },
	Signature("==", TFloat, TFloat):    func(a []any) (any, error) { return a[0].(float64) == a[1].(float64), nil },
	Signature("==", TString, TString):  func(a []any) (any, error) { return a[0].(string) == a[1].(string), nil },
	Signature("==", TBool, TBool):      func(a []any) (any, error) { return a[0].(bool) == a[1].(bool), nil },
	Signature("!=", TInt, TInt):        func(a []any) (any, error) { return a[0].(int) != a[1].(int), nil },
	Signature("!=", TFloat, TFloat):    func(a []any) (any, error) { return a[0].(float64) != a[1].(float64), nil },
	Signature("!=", TString, TString):  func(a []any) (any, error) { return a[0].(string) != a[1].(string), nil },
	Signature("!=", TBool, TBool):      func(a []any) (any, error) { return a[0].(bool) != a[1].(bool), nil },
	Signature(">", TInt, TInt):         func(a []any) (any, error) { return a[0].(int) > a[1].(int), nil },
	Signature(">", TFloat, TFloat):     func(a []any) (any, error) { return a[0].(float64) > a[1].(float64), nil },
	Signature(">", TString, TString):   func(a []any) (any, error) { return a[0].(string) > a[1].(string), nil },
	Signature(">=", TInt, TInt):        func(a []any) (any, error) { return a[0].(int) >= a[1].(int), nil },
	Signature(">=", TFloat, TFloat):    func(a []any) (any, error) { return a[0].(float64) >= a[1].(float64), nil },
	Signature(">=", TString, TString):  func(a []any) (any, error) { return a[0].(string) >= a[1].(string), nil },
	Signature("<", TInt, TInt):         func(a []any) (any, error) { return a[0].(int) < a[1].(int), nil },
	Signature("<", TFloat, TFloat):     func(a []any) (any, error) { return a[0].(float64) < a[1].(float64), nil },
	Signature("<", TString, TString):   func(a []any) (any, error) { return a[0].(string) < a[1].(string), nil },
	Signature("<=", TInt, TInt):        func(a []any) (any, error) { return a[0].(int) <= a[1].(int), nil },
	Signature("<=", TFloat, TFloat):    func(a []any) (any, error) { return a[0].(float64) <= a[1].(float64), nil },
	Signature("<=", TString, TString):  func(a []any) (any, error) { return a[0].(string) <= a[1].(string), nil },
	Signature("+", TInt, TInt):         func(a []any) (any, error) { return a[0].(int) + a[1].(int), nil },
	Signature("+", TFloat, TFloat):     func(a []any) (any, error) { return a[0].(float64) + a[1].(float64), nil },
	Signature("+", TString, TString):   func(a []any) (any, error) { return a[0].(string) + a[1].(string), nil },
	Signature("-", TInt, TInt):         func(a []any) (any, error) { return a[0].(int) - a[1].(int), nil },
	Signature("-", TFloat, TFloat):     func(a []any) (any, error) { return a[0].(float64) - a[1].(float64), nil },
	Signature("-", TInt):               func(a []any) (any, error) { return -a[0].(float64), nil },
	Signature("*", TInt, TInt):         func(a []any) (any, error) { return a[0].(int) * a[1].(int), nil },
	Signature("*", TFloat, TFloat):     func(a []any) (any, error) { return a[0].(float64) * a[1].(float64), nil },
	Signature("/", TInt, TInt):         func(a []any) (any, error) { return a[0].(int) / a[1].(int), nil },
	Signature("/", TFloat, TFloat):     func(a []any) (any, error) { return a[0].(float64) / a[1].(float64), nil },
	Signature("&&", TBool, TBool):      func(a []any) (any, error) { return a[0].(bool) && a[1].(bool), nil },
	Signature("||", TBool, TBool):      func(a []any) (any, error) { return a[0].(bool) || a[1].(bool), nil },
	Signature("!", TBool):              func(a []any) (any, error) { return !a[0].(bool), nil },
	Signature("base64encode", TString): func(a []any) (any, error) { return base64.StdEncoding.EncodeToString([]byte(a[0].(string))), nil },
	Signature("base64decode", TString): func(a []any) (any, error) {
		dec, err := base64.StdEncoding.DecodeString(a[0].(string))
		if err != nil {
			return nil, err
		}
		return string(dec), nil
	},
}

package evaluator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	assert.Equal(t, &Ref{Expr: &Null{}, Key: &String{Value: "a"}}, Parse("a"))
	assert.Equal(t,
		&Ref{
			Expr: &Ref{
				Expr: &Ref{
					Expr: &Null{},
					Key:  &String{Value: "a"},
				},
				Key: &String{Value: "b"},
			},
			Key: &String{Value: "c"},
		},
		Parse("a.b.c"))

}

func TestEvaluate(t *testing.T) {
	expr := Parse("data.config.foo")
	v := Evaluate(expr, map[string]any{
		"data": map[string]any{
			"config": map[string]any{
				"foo": "OK",
			},
		},
	})
	assert.Equal(t, "OK", v)

	expr = Parse("1 + 2 + 4")
	v = Evaluate(expr, map[string]any{})
	assert.Equal(t, 7.0, v)

	expr = Parse(`data.config.foo == "OK"`)
	v = Evaluate(expr, map[string]any{
		"data": map[string]any{
			"config": map[string]any{
				"foo": "OK",
			},
		},
	})
	assert.Equal(t, true, v)

	expr = Parse(`{"low": "Basic", "high": "Premium" }`)
	v = Evaluate(expr, map[string]any{})
	assert.Equal(t, map[Expr]any{
		"low":  "Basic",
		"high": "Premium",
	}, v)

	expr = Parse(`{"low": "Basic", "high": "Premium" }["high"]`)
	v = Evaluate(expr, map[string]any{})
	assert.Equal(t, "Premium", v)
}

func TestCollectReferences(t *testing.T) {
	expr := Parse("modules.foo.xyz + modules.bar.xyz")
	refs := CollectReferences(expr, "^modules\\..*")
	assert.Equal(t, []string{"modules.foo.xyz", "modules.bar.xyz"}, refs)
}

func TestExprArrayReference(t *testing.T) {
	expr := Parse("[1, 2, 3][1]")
	v := Evaluate(expr, map[string]any{})
	assert.Equal(t, 2.0, v)
}

func TestDataArrayReference(t *testing.T) {
	expr := Parse(`list[0].a`)
	v := Evaluate(expr, map[string]any{
		"list": []any{
			map[string]any{
				"a": "foo",
				"b": "bar",
				"c": map[string]any{
					"b": 1.0,
				},
			},
			map[string]any{
				"c": 1.0,
				"d": 2.0,
			},
		},
	})
	assert.Equal(t, 1.0, v)
}

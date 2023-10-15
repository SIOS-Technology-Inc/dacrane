package core

import (
	"github.com/macrat/simplexer"
)

type Module struct {
	Name         string       `yaml:"name"`
	Import       []string     `yaml:"import"`
	Parameter    any          `yaml:"parameter"`
	Dependencies []Dependency `yaml:"dependencies"`
	ModuleCalls  []ModuleCall `yaml:"modules"`
}

type Dependency struct {
	Name   string `yaml:"name"`
	Module string `yaml:"module"`
}

type ModuleCall struct {
	Name      string   `yaml:"name"`
	DependsOn []string `yaml:"depends_on"`
	Module    string   `yaml:"module"`
	Argument  any      `yaml:"argument"`
}

// Expr represents an expression in the AST.
type Expr interface{}

// Number represents a numeric value.
type Number struct {
	Value float64
}

// String represents a string value.
type String struct {
	Value string
}

// Boolean represents a boolean value.
type Boolean struct {
	Value bool
}

// Null represents a null value.
type Null struct{}

// BinaryExpr represents a binary operation.
type BinaryExpr struct {
	Left  Expr
	Op    *simplexer.Token
	Right Expr
}

// UnaryExpr represents a unary operation.
type UnaryExpr struct {
	Op   *simplexer.Token
	Expr Expr
}

// IfExpr represents an if-then-else expression.
type IfExpr struct {
	Condition Expr
	Then      Expr
	Else      Expr
}

// List represents a list of expressions.
type List struct {
	Items []Expr
}

// Map represents a map of string to expression.
type Map struct {
	KVs map[Expr]Expr
}

// App represents a function application.
type App struct {
	Name   string
	Params []Expr
}

// Ref represents a reference to another expression.
type Ref struct {
	Expr Expr
	Key  Expr
}

// Identifier represents an identifier.
type Identifier struct {
	Name string
}

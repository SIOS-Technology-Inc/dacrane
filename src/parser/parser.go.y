%{
package parser

import "github.com/macrat/simplexer"
import "strconv"
import "strings"
import "github.com/SIOS-Technology-Inc/dacrane/v0/src/ast"
import "github.com/SIOS-Technology-Inc/dacrane/v0/src/locator"

func Parse(tokens []*simplexer.Token) (ast.Module, error) {
	lexer := NewTokenIterationLexer(tokens)
	yyParse(lexer)
	if lexer.error != nil {
		return ast.Module{}, lexer.error
	}
	return lexer.result, nil
}
%}

%union{
	token  *simplexer.Token
	Module ast.Module
	Expr   ast.Expr
	Var    ast.Var
	Vars   []ast.Var
}

%left ADD
%left COMMA
%right ASSIGN

%token <token> INTEGER STRING
%token <token> IDENTIFIER
%token <token> ADD
%token <token> LBRACKET RBRACKET

%type <Module> Root Module
%type <Vars> Vars
%type <Var> Var
%type <Expr> Expr
%%

Root: Module {
	yylex.(*TokenIterationLexer).result = $1
	$$ = $1
}

Module: Vars { $$ = ast.Module{ Vars: $1 } }

Vars
  : Var Vars { $$ = append([]ast.Var{ $1 }, $2...) }
  | { $$ = []ast.Var{} }

Var
  : IDENTIFIER ASSIGN Expr {
	  $$ = ast.Var{
			Name: $1.Literal,
			Expr: $3,
			Range: locator.NewRangeFromToken(*$1).Union($3.GetRange()),
		}
	}

Expr
	: INTEGER {
		v, err := strconv.Atoi($1.Literal)
		if err != nil {
			panic(err)
		}
		$$ = &ast.SInt{
			Value: v,
			Range: locator.NewRangeFromToken(*$1),
		}
	}
	| STRING {
		$$ = &ast.SString{
			Value: strings.Replace($1.Literal, "\"", "", -1),
			Range: locator.NewRangeFromToken(*$1),
		}
	}
	| Expr ADD Expr {
		$$ = &ast.App{
			Func: $2.Literal,
			Args: []ast.Expr{$1, $3},
			Range: $1.GetRange().Union($3.GetRange()),
		}
	}
	| IDENTIFIER {
		$$ = &ast.Ref{
			Name: $1.Literal,
			Range: locator.NewRangeFromToken(*$1),
		}
	}
	;

%%

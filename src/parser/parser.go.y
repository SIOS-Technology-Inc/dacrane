%{
package parser

import "github.com/macrat/simplexer"
import "strconv"
import "strings"
import "github.com/SIOS-Technology-Inc/dacrane/v0/src/ast"
import "github.com/SIOS-Technology-Inc/dacrane/v0/src/locator"

func Parse(tokens []*simplexer.Token) ast.Expr {
	lexer := NewTokenIterationLexer(tokens)
	yyParse(lexer)
	return lexer.result
}
%}

%union{
	token *simplexer.Token
	Expr  ast.Expr
	Exprs []ast.Expr
}

%left ADD
%left COMMA

%token <token> INTEGER STRING
%token <token> IDENTIFIER
%token <token> ADD
%token <token> LBRACKET RBRACKET

%type <Expr> Root Expr
%%

Root: Expr {
	yylex.(*TokenIterationLexer).result = $1
	$$ = $1
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
	;

%%

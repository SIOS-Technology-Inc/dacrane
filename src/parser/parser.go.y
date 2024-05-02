%{
package parser

import "github.com/macrat/simplexer"
import "strconv"
import "strings"
import "github.com/SIOS-Technology-Inc/dacrane/v0/src/ast"

func Parse(exprStr string) ast.Expr {
	lexer := NewLexer(strings.NewReader(exprStr))
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
	yylex.(*Lexer).result = $1
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
			Pos: ast.PosFrom(&$1.Position),
		}
	}
	| STRING {
		$$ = &ast.SString{
			Value: strings.Replace($1.Literal, "\"", "", -1),
			Pos: ast.PosFrom(&$1.Position),
		}
	}
	| Expr ADD Expr               { $$ = &ast.App{Func: $2.Literal, Args: []ast.Expr{$1, $3}, Pos: $1.Position()} }
	;

%%

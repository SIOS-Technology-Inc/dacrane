%{
package code

import "github.com/macrat/simplexer"
%}

%union{
	token        *simplexer.Token
	expr         ExprParam
	ref          RefParam
	path         Path
	ident        Identifier
}

%left DOT
%type<expr> expr
%type<ref> ref
%type<path> path
%token<token> IDENT
%token DOT

%%

expr: ref
    {
	    $$ = NewExprParam($1)
		yylex.(*Lexer).result = $$
    }

ref: path { $$ = NewRefParam($1) }

path
	: path DOT path { $$ = append($1, $3...) }
	| IDENT         { $$ = []string{ $1.Literal } }

%%

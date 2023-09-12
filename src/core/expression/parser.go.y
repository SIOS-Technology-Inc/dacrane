%{
package expression

import "github.com/macrat/simplexer"
%}

%union{
	token *simplexer.Token
	path  Path
	ident Identifier
}

%left DOT
%type<path> program path
%token<token> IDENT
%token DOT LBRACKET_EXPR RBRACKET_EXPR

%%

program
	: path
	{
		$$ = $1
		yylex.(*Lexer).result = $1
	}

path
	: path DOT path
	{
		$$ = append($1, $3...)
	}
	| IDENT
	{
		$$ = []string{ $1.Literal }
	}

%%

%{
package parser

import "github.com/macrat/simplexer"
import "strconv"
import "strings"
import "github.com/SIOS-Technology-Inc/dacrane/v0/src/code/ast"

func Parse(exprStr string) ast.Expr {
	lexer := NewLexer(strings.NewReader(exprStr))
	yyParse(lexer)
	return lexer.result
}
%}

%union{
	token    *simplexer.Token
	expr     ast.Expr
	exprs    []ast.Expr
	kvMap    map[ast.Expr]ast.Expr
}

%right IF THEN ELSE
%left LSBRACKET RSBRACKET
%left COMMA
%left OR
%left AND
%left EQ NEQ LT LTE GT GTE
%right NOT
%left ADD SUB
%left MUL DIV
%left PRIORITY
%left DOT
%right UMINUS

%token <token> FLOAT INTEGER
%token <token> STRING
%token <token> BOOLEAN NULL
%token <token> IDENTIFIER
%token <token> DOT COMMA COLON
%token <token> AND OR NOT EQ NEQ LT LTE GT GTE PRIORITY
%token <token> ADD SUB MUL DIV
%token <token> LBRACKET RBRACKET LSBRACKET RSBRACKET LCBRACKET RCBRACKET
%token <token> IF THEN ELSE

%type <expr> Root Expr Ref App Seq Map Variable
%type <exprs> Params Items
%type <kvMap> KVs

%%

Root: Expr {
	yylex.(*Lexer).result = $1
	$$ = $1
}

Expr
	: FLOAT {
		v, err := strconv.ParseFloat($1.Literal, 64)
		if err != nil {
			panic(err)
		}
		$$ = &ast.SFloat{
			Value: v,
			Pos: ast.PosFrom(&$1.Position),
		}
	}
	| INTEGER {
		v, err := strconv.ParseInt($1.Literal, 10, 64)
		if err != nil {
			panic(err)
		}
		$$ = &ast.SInteger{
			Value: v,
			Pos: ast.PosFrom(&$1.Position),
		}
	}
	| STRING  {
		$$ = &ast.SString{
			Value: strings.Replace($1.Literal, "\"", "", -1),
			Pos: ast.PosFrom(&$1.Position),
		}
	}
	| BOOLEAN {
		v, err := strconv.ParseBool($1.Literal)
		if err != nil {
			panic(err)
		}
		$$ = &ast.SBool{
			Value: v,
			Pos: ast.PosFrom(&$1.Position),
		}
	}
	| NULL                        { $$ = &ast.SNull{Pos: ast.PosFrom(&$1.Position)} }
	| LBRACKET Expr RBRACKET      { $$ = $2 }
	| IF Expr THEN Expr ELSE Expr { $$ = &ast.App{Func: "if", Args: []ast.Expr{$4, $6}} }
	| Seq                         { $$ = $1 }
	| Map                         { $$ = $1 }
	| App                         { $$ = $1 }
	| Variable                    { $$ = $1 }
	| Ref                         { $$ = $1 }
	| Expr EQ Expr                { $$ = &ast.App{Func: $2.Literal, Args: []ast.Expr{$1, $3}, Pos: $1.Position()} }
	| Expr NEQ Expr               { $$ = &ast.App{Func: $2.Literal, Args: []ast.Expr{$1, $3}, Pos: $1.Position()} }
	| Expr PRIORITY Expr          { $$ = &ast.App{Func: $2.Literal, Args: []ast.Expr{$1, $3}, Pos: $1.Position()} }
	| Expr LT Expr                { $$ = &ast.App{Func: $2.Literal, Args: []ast.Expr{$1, $3}, Pos: $1.Position()} }
	| Expr LTE Expr               { $$ = &ast.App{Func: $2.Literal, Args: []ast.Expr{$1, $3}, Pos: $1.Position()} }
	| Expr GT Expr                { $$ = &ast.App{Func: $2.Literal, Args: []ast.Expr{$1, $3}, Pos: $1.Position()} }
	| Expr GTE Expr               { $$ = &ast.App{Func: $2.Literal, Args: []ast.Expr{$1, $3}, Pos: $1.Position()} }
	| Expr ADD Expr               { $$ = &ast.App{Func: $2.Literal, Args: []ast.Expr{$1, $3}, Pos: $1.Position()} }
	| Expr SUB Expr               { $$ = &ast.App{Func: $2.Literal, Args: []ast.Expr{$1, $3}, Pos: $1.Position()} }
	| Expr MUL Expr               { $$ = &ast.App{Func: $2.Literal, Args: []ast.Expr{$1, $3}, Pos: $1.Position()} }
	| Expr DIV Expr               { $$ = &ast.App{Func: $2.Literal, Args: []ast.Expr{$1, $3}, Pos: $1.Position()} }
	| Expr AND Expr               { $$ = &ast.App{Func: $2.Literal, Args: []ast.Expr{$1, $3}, Pos: $1.Position()} }
	| Expr OR Expr                { $$ = &ast.App{Func: $2.Literal, Args: []ast.Expr{$1, $3}, Pos: $1.Position()} }
	| NOT Expr                    { $$ = &ast.App{Func: $1.Literal, Args: []ast.Expr{$2}, Pos: ast.PosFrom(&$1.Position)} }
	| SUB Expr %prec UMINUS       { $$ = &ast.App{Func: $1.Literal, Args: []ast.Expr{$2}, Pos: ast.PosFrom(&$1.Position)} }
	;

Ref
	: Expr LSBRACKET Expr RSBRACKET { $$ = &ast.Ref{Dict: $1, Key: $3, Pos: $1.Position()} }
	| Ref DOT IDENTIFIER { $$ = &ast.Ref{Dict: $1, Key: &ast.SString{Value: $3.Literal, Pos: $1.Position()} } }

Variable: IDENTIFIER { $$ = &ast.Variable{Name: $1.Literal, Pos: ast.PosFrom(&$1.Position)} }

App: IDENTIFIER LBRACKET Params RBRACKET { $$ = &ast.App{Func: $1.Literal, Args: $3, Pos: ast.PosFrom(&$1.Position)} }
Params
	: Params COMMA Params { $$ = append($1, $3...) }
	| Expr                { $$ = []ast.Expr{$1} }
	|                     { $$ = []ast.Expr{} }

Seq: LSBRACKET Items RSBRACKET { $$ = &ast.Seq{Value: $2, Pos: ast.PosFrom(&$1.Position)} }
Items
	: Items COMMA Items { $$ = append($1, $3...) }
	| Expr              { $$ = []ast.Expr{$1} }
	|                   { $$ = []ast.Expr{} }

Map: LCBRACKET KVs RCBRACKET { $$ = &ast.Map{Value: $2, Pos: ast.PosFrom(&$1.Position)} }
KVs: KVs COMMA KVs {
		$$ = $1
		for k, v := range $3 {
			$$[k] = v
		}
	}
	| Expr COLON Expr { $$ = map[ast.Expr]ast.Expr{$1: $3} }
	| { $$ = map[ast.Expr]ast.Expr{} }
%%

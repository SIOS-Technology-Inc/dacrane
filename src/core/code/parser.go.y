%{
package code

import "github.com/macrat/simplexer"
import "strconv"
import "strings"
%}

%union{
	token    *simplexer.Token
	expr     Expr
	exprs    []Expr
	kvMap    map[string]Expr
}

%right IF THEN ELSE
%left LSBRACKET RSBRACKET
%left COMMA
%left OR
%left AND
%left EQ LT LTE GT GTE
%right NOT
%left ADD SUB
%left MUL DIV
%left PRIORITY
%left DOT
%right UMINUS

%token <token> NUMBER
%token <token> STRING
%token <token> BOOLEAN NULL
%token <token> IDENTIFIER
%token <token> DOT COMMA COLON
%token <token> AND OR NOT EQ LT LTE GT GTE PRIORITY
%token <token> ADD SUB MUL DIV
%token <token> LBRACKET RBRACKET LSBRACKET RSBRACKET LCBRACKET RCBRACKET
%token <token> IF THEN ELSE

%type <expr> Root Expr Ref App List Map
%type <exprs> Params Items
%type <kvMap> KVs

%%

Root: Expr {
	yylex.(*Lexer).result = $1
	$$ = $1
}

Expr
	: NUMBER {
		v, err := strconv.ParseFloat($1.Literal, 64)
		if (err != nil) {
			panic(err)
		}
		$$ = &Number{Value: v}
	}
	| STRING                      { $$ = &String{Value: strings.Replace($1.Literal, "\"", "", -1)} }
	| BOOLEAN {
		v, err := strconv.ParseBool($1.Literal)
		if (err != nil) {
			panic(err)
		}
		$$ = &Boolean{Value: v}
	}
	| NULL                        { $$ = &Null{} }
	| LBRACKET Expr RBRACKET      { $$ = $2 }
	| IF Expr THEN Expr ELSE Expr { $$ = &IfExpr{Condition: $2, Then: $4, Else: $6} }
	| List                        { $$ = $1 }
	| Map                         { $$ = $1 }
	| App                         { $$ = $1 }
	| Ref                         { $$ = $1 }
	| Expr EQ Expr                { $$ = &BinaryExpr{Left: $1, Op: $2, Right: $3} }
	| Expr LT Expr                { $$ = &BinaryExpr{Left: $1, Op: $2, Right: $3} }
	| Expr LTE Expr               { $$ = &BinaryExpr{Left: $1, Op: $2, Right: $3} }
	| Expr GT Expr                { $$ = &BinaryExpr{Left: $1, Op: $2, Right: $3} }
	| Expr GTE Expr               { $$ = &BinaryExpr{Left: $1, Op: $2, Right: $3} }
	| Expr ADD Expr               { $$ = &BinaryExpr{Left: $1, Op: $2, Right: $3} }
	| Expr SUB Expr               { $$ = &BinaryExpr{Left: $1, Op: $2, Right: $3} }
	| Expr MUL Expr               { $$ = &BinaryExpr{Left: $1, Op: $2, Right: $3} }
	| Expr DIV Expr               { $$ = &BinaryExpr{Left: $1, Op: $2, Right: $3} }
	| Expr AND Expr               { $$ = &BinaryExpr{Left: $1, Op: $2, Right: $3} }
	| Expr OR Expr                { $$ = &BinaryExpr{Left: $1, Op: $2, Right: $3} }
	| NOT Expr                    { $$ = &UnaryExpr{Op: $1, Expr: $2} }
	| SUB Expr %prec UMINUS       { $$ = &UnaryExpr{Op: $1, Expr: $2} }
	;

Ref
	: Expr LSBRACKET Expr RSBRACKET { $$ = &Ref{Expr: $1, Key: $3} }
	| IDENTIFIER { $$ = &Identifier{Name: $1.Literal} }
	| Ref DOT IDENTIFIER { $$ = &Identifier{Name: $1.(*Identifier).Name + "." + $3.Literal} }


App: IDENTIFIER LBRACKET Params RBRACKET { $$ = &App{Name: $1.Literal, Params: $3} }
Params
	: Params COMMA Params { $$ = append($1, $3...) }
	| Expr                { $$ = []Expr{$1} }
	|                     { $$ = []Expr{} }

List: LSBRACKET Items RSBRACKET { $$ = &List{Items: $2} }
Items
	: Items COMMA Items { $$ = append($1, $3...) }
	| Expr              { $$ = []Expr{$1} }
	|                   { $$ = []Expr{} }

Map: LCBRACKET KVs RCBRACKET { $$ = &Map{KVs: $2} }
KVs: KVs COMMA KVs {
		$$ = $1
		for k, v := range $3 {
			$$[k] = v
		}
	}
	| STRING COLON Expr { $$ = map[string]Expr{$1.Literal: $3} }
	| { $$ = map[string]Expr{} }
%%

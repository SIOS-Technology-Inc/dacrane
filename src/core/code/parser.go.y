%{
package code

import "github.com/macrat/simplexer"
%}

%union{
	token    *simplexer.Token
	tnumber  float32
	tboolean bool
  tstring  string
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

%token<tnumber> NUMBER
%token<tstring> STRING
%token<tboolean> BOOLEAN
%token NULL
%token<tstring> IDENTIFIER
%token DOT COMMA COLON
%token AND OR NOT EQ LT LTE GT GTE PRIORITY
%token ADD SUB MUL DIV
%token LBRACKET RBRACKET LSBRACKET RSBRACKET LCBRACKET RCBRACKET
%token IF THEN ELSE

%%

Expr
	: NUMBER
	| STRING
	| BOOLEAN
	| NULL
	| LBRACKET Expr RBRACKET
	| IF Expr THEN Expr ELSE Expr
	| List
	| Map
	| App
	| Ref
	| Expr EQ Expr
	| Expr LT Expr
	| Expr LTE Expr
	| Expr GT Expr
	| Expr GTE Expr
	| NOT Expr
	| Expr PRIORITY Expr
	| Expr ADD Expr
	| Expr SUB Expr
	| Expr MUL Expr
	| Expr DIV Expr
	| SUB Expr %prec UMINUS
	;

Ref
  : Expr LSBRACKET Expr RSBRACKET
	| IDENTIFIER DOT IDENTIFIER

App: IDENTIFIER LBRACKET Params RBRACKET
Params
  : Params COMMA Params
	| Expr
	|

List: LSBRACKET Items RSBRACKET
Items
  : Items COMMA Items
	| Expr
	|

Map: LCBRACKET KVs RCBRACKET
KVs
  : KVs COMMA KVs
	| STRING COLON Expr
	|
%%

// Code generated by goyacc -o core/code/parser.go core/code/parser.go.y. DO NOT EDIT.

//line core/code/parser.go.y:2
package core

import __yyfmt__ "fmt"

//line core/code/parser.go.y:2

import "github.com/macrat/simplexer"
import "strconv"
import "strings"

//line core/code/parser.go.y:9
type yySymType struct {
	yys   int
	token *simplexer.Token
	expr  Expr
	exprs []Expr
	kvMap map[Expr]Expr
}

const IF = 57346
const THEN = 57347
const ELSE = 57348
const LSBRACKET = 57349
const RSBRACKET = 57350
const COMMA = 57351
const OR = 57352
const AND = 57353
const EQ = 57354
const LT = 57355
const LTE = 57356
const GT = 57357
const GTE = 57358
const NOT = 57359
const ADD = 57360
const SUB = 57361
const MUL = 57362
const DIV = 57363
const PRIORITY = 57364
const DOT = 57365
const UMINUS = 57366
const NUMBER = 57367
const STRING = 57368
const BOOLEAN = 57369
const NULL = 57370
const IDENTIFIER = 57371
const COLON = 57372
const LBRACKET = 57373
const RBRACKET = 57374
const LCBRACKET = 57375
const RCBRACKET = 57376

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"IF",
	"THEN",
	"ELSE",
	"LSBRACKET",
	"RSBRACKET",
	"COMMA",
	"OR",
	"AND",
	"EQ",
	"LT",
	"LTE",
	"GT",
	"GTE",
	"NOT",
	"ADD",
	"SUB",
	"MUL",
	"DIV",
	"PRIORITY",
	"DOT",
	"UMINUS",
	"NUMBER",
	"STRING",
	"BOOLEAN",
	"NULL",
	"IDENTIFIER",
	"COLON",
	"LBRACKET",
	"RBRACKET",
	"LCBRACKET",
	"RCBRACKET",
}

var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line core/code/parser.go.y:116

//line yacctab:1
var yyExca = [...]int8{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyPrivate = 57344

const yyLast = 204

var yyAct = [...]int8{
	37, 2, 62, 38, 36, 40, 56, 60, 31, 32,
	70, 26, 27, 19, 34, 35, 33, 39, 19, 41,
	42, 43, 44, 45, 46, 47, 48, 49, 50, 51,
	52, 53, 59, 69, 24, 25, 26, 27, 19, 10,
	30, 63, 9, 29, 28, 18, 20, 21, 22, 23,
	11, 24, 25, 26, 27, 19, 65, 57, 58, 12,
	1, 39, 68, 66, 67, 54, 8, 0, 0, 15,
	0, 63, 73, 72, 0, 0, 0, 0, 0, 13,
	0, 14, 0, 0, 0, 0, 0, 3, 4, 5,
	6, 17, 0, 7, 30, 16, 0, 29, 28, 18,
	20, 21, 22, 23, 0, 24, 25, 26, 27, 19,
	0, 0, 0, 0, 71, 30, 0, 61, 29, 28,
	18, 20, 21, 22, 23, 0, 24, 25, 26, 27,
	19, 30, 64, 0, 29, 28, 18, 20, 21, 22,
	23, 0, 24, 25, 26, 27, 19, 55, 0, 30,
	0, 0, 29, 28, 18, 20, 21, 22, 23, 0,
	24, 25, 26, 27, 19, 30, 0, 0, 29, 28,
	18, 20, 21, 22, 23, 0, 24, 25, 26, 27,
	19, 28, 18, 20, 21, 22, 23, 0, 24, 25,
	26, 27, 19, 18, 20, 21, 22, 23, 0, 24,
	25, 26, 27, 19,
}

var yyPact = [...]int16{
	62, -1000, 158, -1000, -1000, -1000, -1000, 62, 62, -1000,
	-1000, -1000, -7, 62, 62, 62, 62, -26, 62, 62,
	62, 62, 62, 62, 62, 62, 62, 62, 62, 62,
	62, 33, 142, -23, 16, -1000, 49, 158, -2, 87,
	62, 16, -1000, 16, 16, 16, 16, -9, -9, -4,
	-4, 181, 170, 124, -1000, 62, -1000, -1000, 62, -1000,
	62, 62, 1, 158, -1000, 108, -1000, -1000, 158, -1000,
	62, 62, -1000, 158,
}

var yyPgo = [...]int8{
	0, 60, 0, 59, 50, 42, 39, 2, 4, 3,
}

var yyR1 = [...]int8{
	0, 1, 2, 2, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 2, 2, 3, 3, 3, 4,
	7, 7, 7, 5, 8, 8, 8, 6, 9, 9,
	9,
}

var yyR2 = [...]int8{
	0, 1, 1, 1, 1, 1, 3, 6, 1, 1,
	1, 1, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 2, 2, 4, 1, 3, 4,
	3, 1, 0, 3, 3, 1, 0, 3, 3, 3,
	0,
}

var yyChk = [...]int16{
	-1000, -1, -2, 25, 26, 27, 28, 31, 4, -5,
	-6, -4, -3, 17, 19, 7, 33, 29, 12, 22,
	13, 14, 15, 16, 18, 19, 20, 21, 11, 10,
	7, -2, -2, 23, -2, -2, -8, -2, -9, -2,
	31, -2, -2, -2, -2, -2, -2, -2, -2, -2,
	-2, -2, -2, -2, 32, 5, 29, 8, 9, 34,
	9, 30, -7, -2, 8, -2, -8, -9, -2, 32,
	9, 6, -7, -2,
}

var yyDef = [...]int8{
	0, -2, 1, 2, 3, 4, 5, 0, 0, 8,
	9, 10, 11, 0, 0, 36, 40, 27, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 24, 25, 0, 35, 0, 0,
	32, 12, 13, 14, 15, 16, 17, 18, 19, 20,
	21, 22, 23, 0, 6, 0, 28, 33, 36, 37,
	40, 0, 0, 31, 26, 0, 34, 38, 39, 29,
	32, 0, 30, 7,
}

var yyTok1 = [...]int8{
	1,
}

var yyTok2 = [...]int8{
	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34,
}

var yyTok3 = [...]int8{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := int(yyPact[state])
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && int(yyChk[int(yyAct[n])]) == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || int(yyExca[i+1]) != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := int(yyExca[i])
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = int(yyTok1[0])
		goto out
	}
	if char < len(yyTok1) {
		token = int(yyTok1[char])
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = int(yyTok2[char-yyPrivate])
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = int(yyTok3[i+0])
		if token == char {
			token = int(yyTok3[i+1])
			goto out
		}
	}

out:
	if token == 0 {
		token = int(yyTok2[1]) /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = int(yyPact[yystate])
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = int(yyAct[yyn])
	if int(yyChk[yyn]) == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = int(yyDef[yystate])
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && int(yyExca[xi+1]) == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = int(yyExca[xi+0])
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = int(yyExca[xi+1])
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = int(yyPact[yyS[yyp].yys]) + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = int(yyAct[yyn]) /* simulate a shift of "error" */
					if int(yyChk[yystate]) == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= int(yyR2[yyn])
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = int(yyR1[yyn])
	yyg := int(yyPgo[yyn])
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = int(yyAct[yyg])
	} else {
		yystate = int(yyAct[yyj])
		if int(yyChk[yystate]) != -yyn {
			yystate = int(yyAct[yyg])
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-1 : yypt+1]
//line core/code/parser.go.y:45
		{
			yylex.(*Lexer).result = yyDollar[1].expr
			yyVAL.expr = yyDollar[1].expr
		}
	case 2:
		yyDollar = yyS[yypt-1 : yypt+1]
//line core/code/parser.go.y:51
		{
			v, err := strconv.ParseFloat(yyDollar[1].token.Literal, 64)
			if err != nil {
				panic(err)
			}
			yyVAL.expr = &Number{Value: v}
		}
	case 3:
		yyDollar = yyS[yypt-1 : yypt+1]
//line core/code/parser.go.y:58
		{
			yyVAL.expr = &String{Value: strings.Replace(yyDollar[1].token.Literal, "\"", "", -1)}
		}
	case 4:
		yyDollar = yyS[yypt-1 : yypt+1]
//line core/code/parser.go.y:59
		{
			v, err := strconv.ParseBool(yyDollar[1].token.Literal)
			if err != nil {
				panic(err)
			}
			yyVAL.expr = &Boolean{Value: v}
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
//line core/code/parser.go.y:66
		{
			yyVAL.expr = &Null{}
		}
	case 6:
		yyDollar = yyS[yypt-3 : yypt+1]
//line core/code/parser.go.y:67
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 7:
		yyDollar = yyS[yypt-6 : yypt+1]
//line core/code/parser.go.y:68
		{
			yyVAL.expr = &IfExpr{Condition: yyDollar[2].expr, Then: yyDollar[4].expr, Else: yyDollar[6].expr}
		}
	case 8:
		yyDollar = yyS[yypt-1 : yypt+1]
//line core/code/parser.go.y:69
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 9:
		yyDollar = yyS[yypt-1 : yypt+1]
//line core/code/parser.go.y:70
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
//line core/code/parser.go.y:71
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 11:
		yyDollar = yyS[yypt-1 : yypt+1]
//line core/code/parser.go.y:72
		{
			yyVAL.expr = yyDollar[1].expr
		}
	case 12:
		yyDollar = yyS[yypt-3 : yypt+1]
//line core/code/parser.go.y:73
		{
			yyVAL.expr = &BinaryExpr{Left: yyDollar[1].expr, Op: yyDollar[2].token, Right: yyDollar[3].expr}
		}
	case 13:
		yyDollar = yyS[yypt-3 : yypt+1]
//line core/code/parser.go.y:74
		{
			yyVAL.expr = &BinaryExpr{Left: yyDollar[1].expr, Op: yyDollar[2].token, Right: yyDollar[3].expr}
		}
	case 14:
		yyDollar = yyS[yypt-3 : yypt+1]
//line core/code/parser.go.y:75
		{
			yyVAL.expr = &BinaryExpr{Left: yyDollar[1].expr, Op: yyDollar[2].token, Right: yyDollar[3].expr}
		}
	case 15:
		yyDollar = yyS[yypt-3 : yypt+1]
//line core/code/parser.go.y:76
		{
			yyVAL.expr = &BinaryExpr{Left: yyDollar[1].expr, Op: yyDollar[2].token, Right: yyDollar[3].expr}
		}
	case 16:
		yyDollar = yyS[yypt-3 : yypt+1]
//line core/code/parser.go.y:77
		{
			yyVAL.expr = &BinaryExpr{Left: yyDollar[1].expr, Op: yyDollar[2].token, Right: yyDollar[3].expr}
		}
	case 17:
		yyDollar = yyS[yypt-3 : yypt+1]
//line core/code/parser.go.y:78
		{
			yyVAL.expr = &BinaryExpr{Left: yyDollar[1].expr, Op: yyDollar[2].token, Right: yyDollar[3].expr}
		}
	case 18:
		yyDollar = yyS[yypt-3 : yypt+1]
//line core/code/parser.go.y:79
		{
			yyVAL.expr = &BinaryExpr{Left: yyDollar[1].expr, Op: yyDollar[2].token, Right: yyDollar[3].expr}
		}
	case 19:
		yyDollar = yyS[yypt-3 : yypt+1]
//line core/code/parser.go.y:80
		{
			yyVAL.expr = &BinaryExpr{Left: yyDollar[1].expr, Op: yyDollar[2].token, Right: yyDollar[3].expr}
		}
	case 20:
		yyDollar = yyS[yypt-3 : yypt+1]
//line core/code/parser.go.y:81
		{
			yyVAL.expr = &BinaryExpr{Left: yyDollar[1].expr, Op: yyDollar[2].token, Right: yyDollar[3].expr}
		}
	case 21:
		yyDollar = yyS[yypt-3 : yypt+1]
//line core/code/parser.go.y:82
		{
			yyVAL.expr = &BinaryExpr{Left: yyDollar[1].expr, Op: yyDollar[2].token, Right: yyDollar[3].expr}
		}
	case 22:
		yyDollar = yyS[yypt-3 : yypt+1]
//line core/code/parser.go.y:83
		{
			yyVAL.expr = &BinaryExpr{Left: yyDollar[1].expr, Op: yyDollar[2].token, Right: yyDollar[3].expr}
		}
	case 23:
		yyDollar = yyS[yypt-3 : yypt+1]
//line core/code/parser.go.y:84
		{
			yyVAL.expr = &BinaryExpr{Left: yyDollar[1].expr, Op: yyDollar[2].token, Right: yyDollar[3].expr}
		}
	case 24:
		yyDollar = yyS[yypt-2 : yypt+1]
//line core/code/parser.go.y:85
		{
			yyVAL.expr = &UnaryExpr{Op: yyDollar[1].token, Expr: yyDollar[2].expr}
		}
	case 25:
		yyDollar = yyS[yypt-2 : yypt+1]
//line core/code/parser.go.y:86
		{
			yyVAL.expr = &UnaryExpr{Op: yyDollar[1].token, Expr: yyDollar[2].expr}
		}
	case 26:
		yyDollar = yyS[yypt-4 : yypt+1]
//line core/code/parser.go.y:90
		{
			yyVAL.expr = &Ref{Expr: yyDollar[1].expr, Key: yyDollar[3].expr}
		}
	case 27:
		yyDollar = yyS[yypt-1 : yypt+1]
//line core/code/parser.go.y:91
		{
			yyVAL.expr = &Identifier{Name: yyDollar[1].token.Literal}
		}
	case 28:
		yyDollar = yyS[yypt-3 : yypt+1]
//line core/code/parser.go.y:92
		{
			yyVAL.expr = &Identifier{Name: yyDollar[1].expr.(*Identifier).Name + "." + yyDollar[3].token.Literal}
		}
	case 29:
		yyDollar = yyS[yypt-4 : yypt+1]
//line core/code/parser.go.y:95
		{
			yyVAL.expr = &App{Name: yyDollar[1].token.Literal, Params: yyDollar[3].exprs}
		}
	case 30:
		yyDollar = yyS[yypt-3 : yypt+1]
//line core/code/parser.go.y:97
		{
			yyVAL.exprs = append(yyDollar[1].exprs, yyDollar[3].exprs...)
		}
	case 31:
		yyDollar = yyS[yypt-1 : yypt+1]
//line core/code/parser.go.y:98
		{
			yyVAL.exprs = []Expr{yyDollar[1].expr}
		}
	case 32:
		yyDollar = yyS[yypt-0 : yypt+1]
//line core/code/parser.go.y:99
		{
			yyVAL.exprs = []Expr{}
		}
	case 33:
		yyDollar = yyS[yypt-3 : yypt+1]
//line core/code/parser.go.y:101
		{
			yyVAL.expr = &List{Items: yyDollar[2].exprs}
		}
	case 34:
		yyDollar = yyS[yypt-3 : yypt+1]
//line core/code/parser.go.y:103
		{
			yyVAL.exprs = append(yyDollar[1].exprs, yyDollar[3].exprs...)
		}
	case 35:
		yyDollar = yyS[yypt-1 : yypt+1]
//line core/code/parser.go.y:104
		{
			yyVAL.exprs = []Expr{yyDollar[1].expr}
		}
	case 36:
		yyDollar = yyS[yypt-0 : yypt+1]
//line core/code/parser.go.y:105
		{
			yyVAL.exprs = []Expr{}
		}
	case 37:
		yyDollar = yyS[yypt-3 : yypt+1]
//line core/code/parser.go.y:107
		{
			yyVAL.expr = &Map{KVs: yyDollar[2].kvMap}
		}
	case 38:
		yyDollar = yyS[yypt-3 : yypt+1]
//line core/code/parser.go.y:108
		{
			yyVAL.kvMap = yyDollar[1].kvMap
			for k, v := range yyDollar[3].kvMap {
				yyVAL.kvMap[k] = v
			}
		}
	case 39:
		yyDollar = yyS[yypt-3 : yypt+1]
//line core/code/parser.go.y:114
		{
			yyVAL.kvMap = map[Expr]Expr{yyDollar[1].expr: yyDollar[3].expr}
		}
	case 40:
		yyDollar = yyS[yypt-0 : yypt+1]
//line core/code/parser.go.y:115
		{
			yyVAL.kvMap = map[Expr]Expr{}
		}
	}
	goto yystack /* stack new state and value */
}
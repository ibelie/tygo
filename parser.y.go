//line parser.y:10
package tygo

import __yyfmt__ "fmt"

//line parser.y:11
import (
	"bytes"
	"log"
	"strconv"
	"unicode/utf8"
)

//line parser.y:22
type tygoSymType struct {
	yys     int
	ident   string
	keyword string
	integer int
	spec    Type
	specs   []Type
	method  *Method
	object  *Object
	enum    *Enum
}

const TYPE = 57346
const ENUM = 57347
const OBJECT = 57348
const MAP = 57349
const FIXEDPOINT = 57350
const VARIANT = 57351
const IOTA = 57352
const IDENT = 57353
const INTEGER = 57354

var tygoToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"'*'",
	"'='",
	"'.'",
	"','",
	"'('",
	"')'",
	"'['",
	"']'",
	"'{'",
	"'}'",
	"'<'",
	"'>'",
	"'\\t'",
	"'\\n'",
	"TYPE",
	"ENUM",
	"OBJECT",
	"MAP",
	"FIXEDPOINT",
	"VARIANT",
	"IOTA",
	"IDENT",
	"INTEGER",
}
var tygoStatenames = [...]string{}

const tygoEofCode = 1
const tygoErrCode = 2
const tygoInitialStackSize = 16

//line parser.y:231

var eiota int

var (
	parserPkg     string
	parserTypes   []Type
	parserTypeMap map[string]Type
	parserImports map[string]string
	parserTypePkg map[string][2]string
)

func Parse(code string, pkg string, imports map[string]string, typePkg map[string][2]string) []Type {
	parserPkg = pkg
	parserTypes = nil
	parserTypeMap = make(map[string]Type)
	parserImports = imports
	parserTypePkg = typePkg
	tygoParse(&tygoLex{code: []byte(code)})
	return parserTypes
}

// The parser expects the lexer to return 0 on EOF.  Give it a name for clarity.
const eof = 0

// The parser uses the type <prefix>Lex as a lexer. It must provide
// the methods Lex(*<prefix>SymType) int and Error(string).
type tygoLex struct {
	code []byte
	peek rune
}

// The parser calls this method to get each new token. This
// implementation returns symbols and ID.
func (x *tygoLex) Lex(yylval *tygoSymType) int {
	for {
		c := x.next()
		switch c {
		case eof:
			return eof
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return x.integer(c, yylval)
		case '*', '=', '.', ',', '(', ')', '[', ']', '{', '}', '<', '>', '\t', '\n':
			return int(c)
		case ' ', '\r':
		default:
			return x.ident(c, yylval)
		}
	}
}

// Lex a integer.
func (x *tygoLex) integer(c rune, yylval *tygoSymType) int {
	add := func(b *bytes.Buffer, c rune) {
		if _, err := b.WriteRune(c); err != nil {
			log.Fatalf("[Tygo][Parser] WriteRune: %s", err)
		}
	}
	var b bytes.Buffer
	add(&b, c)
L:
	for {
		c = x.next()
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			add(&b, c)
		default:
			break L
		}
	}
	if c != eof {
		x.peek = c
	}
	if i, err := strconv.Atoi(b.String()); err == nil {
		yylval.integer = i
	} else {
		log.Fatalf("[Tygo][Parser] integer: %s", err)
	}
	return INTEGER
}

// Lex a ident/keyword.
func (x *tygoLex) ident(c rune, yylval *tygoSymType) int {
	add := func(b *bytes.Buffer, c rune) {
		if _, err := b.WriteRune(c); err != nil {
			log.Fatalf("[Tygo][Parser] WriteRune: %s", err)
		}
	}
	var b bytes.Buffer
	add(&b, c)
L:
	for {
		c = x.next()
		switch c {
		case '*', '=', '.', ',', '(', ')', '[', ']', '{', '}', '<', '>', '\t', '\n', ' ', '\r':
			break L
		default:
			add(&b, c)
		}
	}
	if c != eof {
		x.peek = c
	}
	yylval.keyword = b.String()
	switch yylval.keyword {
	case "type":
		return TYPE
	case "enum":
		return ENUM
	case "object":
		return OBJECT
	case "map":
		return MAP
	case "fixedpoint":
		return FIXEDPOINT
	case "variant":
		return VARIANT
	case "iota":
		return IOTA
	default:
		yylval.ident = yylval.keyword
		return IDENT
	}
}

// Return the next rune for the lexer.
func (x *tygoLex) next() rune {
	if x.peek != eof {
		r := x.peek
		x.peek = eof
		return r
	}
	if len(x.code) == 0 {
		return eof
	}
	c, size := utf8.DecodeRune(x.code)
	x.code = x.code[size:]
	if c == utf8.RuneError && size == 1 {
		log.Fatalf("[Tygo][Parser] Invalid utf8")
		return x.next()
	}
	return c
}

// The parser calls this method on a parse error.
func (x *tygoLex) Error(s string) {
	log.Fatalf("[Tygo][Parser] Parse error: %s", s)
}

//line yacctab:1
var tygoExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const tygoPrivate = 57344

const tygoLast = 115

var tygoAct = [...]int{

	14, 64, 53, 49, 52, 87, 76, 22, 78, 19,
	17, 55, 43, 27, 28, 41, 16, 31, 11, 4,
	32, 42, 30, 22, 44, 45, 48, 46, 18, 36,
	25, 26, 29, 54, 15, 15, 15, 88, 59, 58,
	37, 39, 38, 60, 40, 13, 62, 68, 10, 51,
	63, 69, 70, 71, 72, 83, 50, 82, 73, 74,
	22, 77, 75, 56, 12, 47, 36, 8, 80, 22,
	81, 34, 35, 33, 57, 36, 22, 37, 39, 38,
	85, 40, 36, 20, 86, 84, 37, 39, 38, 9,
	40, 66, 10, 37, 39, 38, 7, 40, 68, 8,
	79, 68, 66, 67, 65, 61, 34, 3, 2, 6,
	5, 1, 24, 23, 21,
}
var tygoPact = [...]int{

	1, 1, 83, 76, -7, 51, 32, 18, -9, 18,
	3, 11, 18, 18, 15, -1000, 17, 15, 65, -10,
	18, -1000, -13, 19, 56, 44, 37, 15, 15, -1000,
	-22, 15, 18, -1000, -14, -1000, 52, 64, 25, 24,
	100, 72, 15, 99, 15, 18, 72, -1000, 95, 94,
	18, 18, 18, 18, 15, -1000, 72, 72, 72, -20,
	18, -17, 15, 91, 84, -1000, 72, -1000, 72, 15,
	15, 15, 15, -1000, 46, 40, 78, 15, -1000, 18,
	-1000, -1000, 72, -1000, -21, 15, -1000, 22, -1000,
}
var tygoPgo = [...]int{

	0, 108, 107, 114, 113, 112, 3, 1, 72, 111,
	0,
}
var tygoR1 = [...]int{

	0, 9, 9, 9, 9, 1, 1, 1, 1, 2,
	2, 2, 2, 2, 3, 3, 3, 4, 4, 4,
	5, 6, 6, 7, 7, 7, 7, 7, 8, 8,
	8, 8, 10, 10,
}
var tygoR2 = [...]int{

	0, 3, 4, 3, 4, 5, 6, 6, 4, 5,
	5, 6, 4, 3, 2, 3, 5, 2, 3, 3,
	2, 3, 3, 1, 3, 5, 4, 6, 1, 3,
	2, 4, 1, 2,
}
var tygoChk = [...]int{

	-1000, -9, -1, -2, 18, -1, -2, 13, 16, 13,
	16, 25, 13, 13, -10, 17, 25, -10, 25, 6,
	-8, -3, 4, -4, -5, 19, 20, -10, -10, 17,
	5, -10, -7, 8, 6, -8, 10, 21, 23, 22,
	25, 25, -10, 25, -10, -7, 8, 9, -7, -6,
	12, 12, 26, 24, -10, 25, 11, 10, 14, 14,
	-7, 6, -10, -6, -7, 9, 7, 9, 7, -10,
	-10, -10, -10, -7, -7, -6, 26, -10, 25, 9,
	-7, -7, 11, 15, 7, -10, -7, 26, 15,
}
var tygoDef = [...]int{

	0, -2, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 1, 32, 0, 3, 28, 0,
	0, 13, 0, 0, 0, 0, 0, 2, 4, 33,
	0, 8, 0, 20, 0, 23, 0, 0, 0, 0,
	28, 0, 12, 30, 14, 0, 0, 17, 0, 0,
	0, 0, 0, 0, 10, 29, 0, 0, 0, 0,
	0, 0, 15, 0, 0, 18, 0, 19, 0, 5,
	9, 6, 7, 24, 0, 0, 0, 11, 31, 0,
	21, 22, 0, 26, 0, 16, 25, 0, 27,
}
var tygoTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 16,
	17, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	8, 9, 4, 3, 7, 3, 6, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	14, 5, 15, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 10, 3, 11, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 12, 3, 13,
}
var tygoTok2 = [...]int{

	2, 3, 18, 19, 20, 21, 22, 23, 24, 25,
	26,
}
var tygoTok3 = [...]int{
	0,
}

var tygoErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	tygoDebug        = 0
	tygoErrorVerbose = false
)

type tygoLexer interface {
	Lex(lval *tygoSymType) int
	Error(s string)
}

type tygoParser interface {
	Parse(tygoLexer) int
	Lookahead() int
}

type tygoParserImpl struct {
	lval  tygoSymType
	stack [tygoInitialStackSize]tygoSymType
	char  int
}

func (p *tygoParserImpl) Lookahead() int {
	return p.char
}

func tygoNewParser() tygoParser {
	return &tygoParserImpl{}
}

const tygoFlag = -1000

func tygoTokname(c int) string {
	if c >= 1 && c-1 < len(tygoToknames) {
		if tygoToknames[c-1] != "" {
			return tygoToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func tygoStatname(s int) string {
	if s >= 0 && s < len(tygoStatenames) {
		if tygoStatenames[s] != "" {
			return tygoStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func tygoErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !tygoErrorVerbose {
		return "syntax error"
	}

	for _, e := range tygoErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + tygoTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := tygoPact[state]
	for tok := TOKSTART; tok-1 < len(tygoToknames); tok++ {
		if n := base + tok; n >= 0 && n < tygoLast && tygoChk[tygoAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if tygoDef[state] == -2 {
		i := 0
		for tygoExca[i] != -1 || tygoExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; tygoExca[i] >= 0; i += 2 {
			tok := tygoExca[i]
			if tok < TOKSTART || tygoExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if tygoExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += tygoTokname(tok)
	}
	return res
}

func tygolex1(lex tygoLexer, lval *tygoSymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = tygoTok1[0]
		goto out
	}
	if char < len(tygoTok1) {
		token = tygoTok1[char]
		goto out
	}
	if char >= tygoPrivate {
		if char < tygoPrivate+len(tygoTok2) {
			token = tygoTok2[char-tygoPrivate]
			goto out
		}
	}
	for i := 0; i < len(tygoTok3); i += 2 {
		token = tygoTok3[i+0]
		if token == char {
			token = tygoTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = tygoTok2[1] /* unknown char */
	}
	if tygoDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", tygoTokname(token), uint(char))
	}
	return char, token
}

func tygoParse(tygolex tygoLexer) int {
	return tygoNewParser().Parse(tygolex)
}

func (tygorcvr *tygoParserImpl) Parse(tygolex tygoLexer) int {
	var tygon int
	var tygoVAL tygoSymType
	var tygoDollar []tygoSymType
	_ = tygoDollar // silence set and not used
	tygoS := tygorcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	tygostate := 0
	tygorcvr.char = -1
	tygotoken := -1 // tygorcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		tygostate = -1
		tygorcvr.char = -1
		tygotoken = -1
	}()
	tygop := -1
	goto tygostack

ret0:
	return 0

ret1:
	return 1

tygostack:
	/* put a state and value onto the stack */
	if tygoDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", tygoTokname(tygotoken), tygoStatname(tygostate))
	}

	tygop++
	if tygop >= len(tygoS) {
		nyys := make([]tygoSymType, len(tygoS)*2)
		copy(nyys, tygoS)
		tygoS = nyys
	}
	tygoS[tygop] = tygoVAL
	tygoS[tygop].yys = tygostate

tygonewstate:
	tygon = tygoPact[tygostate]
	if tygon <= tygoFlag {
		goto tygodefault /* simple state */
	}
	if tygorcvr.char < 0 {
		tygorcvr.char, tygotoken = tygolex1(tygolex, &tygorcvr.lval)
	}
	tygon += tygotoken
	if tygon < 0 || tygon >= tygoLast {
		goto tygodefault
	}
	tygon = tygoAct[tygon]
	if tygoChk[tygon] == tygotoken { /* valid shift */
		tygorcvr.char = -1
		tygotoken = -1
		tygoVAL = tygorcvr.lval
		tygostate = tygon
		if Errflag > 0 {
			Errflag--
		}
		goto tygostack
	}

tygodefault:
	/* default state action */
	tygon = tygoDef[tygostate]
	if tygon == -2 {
		if tygorcvr.char < 0 {
			tygorcvr.char, tygotoken = tygolex1(tygolex, &tygorcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if tygoExca[xi+0] == -1 && tygoExca[xi+1] == tygostate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			tygon = tygoExca[xi+0]
			if tygon < 0 || tygon == tygotoken {
				break
			}
		}
		tygon = tygoExca[xi+1]
		if tygon < 0 {
			goto ret0
		}
	}
	if tygon == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			tygolex.Error(tygoErrorMessage(tygostate, tygotoken))
			Nerrs++
			if tygoDebug >= 1 {
				__yyfmt__.Printf("%s", tygoStatname(tygostate))
				__yyfmt__.Printf(" saw %s\n", tygoTokname(tygotoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for tygop >= 0 {
				tygon = tygoPact[tygoS[tygop].yys] + tygoErrCode
				if tygon >= 0 && tygon < tygoLast {
					tygostate = tygoAct[tygon] /* simulate a shift of "error" */
					if tygoChk[tygostate] == tygoErrCode {
						goto tygostack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if tygoDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", tygoS[tygop].yys)
				}
				tygop--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if tygoDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", tygoTokname(tygotoken))
			}
			if tygotoken == tygoEofCode {
				goto ret1
			}
			tygorcvr.char = -1
			tygotoken = -1
			goto tygonewstate /* try again in the same state */
		}
	}

	/* reduction by production tygon */
	if tygoDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", tygon, tygoStatname(tygostate))
	}

	tygont := tygon
	tygopt := tygop
	_ = tygopt // guard against "declared and not used"

	tygop -= tygoR2[tygon]
	// tygop is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if tygop+1 >= len(tygoS) {
		nyys := make([]tygoSymType, len(tygoS)*2)
		copy(nyys, tygoS)
		tygoS = nyys
	}
	tygoVAL = tygoS[tygop+1]

	/* consult goto table to find next state */
	tygon = tygoR1[tygon]
	tygog := tygoPgo[tygon]
	tygoj := tygog + tygoS[tygop].yys + 1

	if tygoj >= tygoLast {
		tygostate = tygoAct[tygog]
	} else {
		tygostate = tygoAct[tygoj]
		if tygoChk[tygostate] != -tygon {
			tygostate = tygoAct[tygog]
		}
	}
	// dummy call; replaced with literal code
	switch tygont {

	case 1:
		tygoDollar = tygoS[tygopt-3 : tygopt+1]
		//line parser.y:48
		{
			parserTypes = append(parserTypes, tygoDollar[1].enum)
			parserTypeMap[tygoDollar[1].enum.Name] = tygoDollar[1].enum
		}
	case 2:
		tygoDollar = tygoS[tygopt-4 : tygopt+1]
		//line parser.y:53
		{
			parserTypes = append(parserTypes, tygoDollar[2].enum)
			parserTypeMap[tygoDollar[2].enum.Name] = tygoDollar[2].enum
		}
	case 3:
		tygoDollar = tygoS[tygopt-3 : tygopt+1]
		//line parser.y:58
		{
			if tygoDollar[1].object.Parent == nil {
				tygoDollar[1].object.Parent = &InstanceType{PkgName: "tygo", PkgPath: TYGO_PATH, Name: "Tygo"}
			}
			parserTypes = append(parserTypes, tygoDollar[1].object)
			parserTypeMap[tygoDollar[1].object.Name] = tygoDollar[1].object
		}
	case 4:
		tygoDollar = tygoS[tygopt-4 : tygopt+1]
		//line parser.y:66
		{
			if tygoDollar[2].object.Parent == nil {
				tygoDollar[2].object.Parent = &InstanceType{PkgName: "tygo", PkgPath: TYGO_PATH, Name: "Tygo"}
			}
			parserTypes = append(parserTypes, tygoDollar[2].object)
			parserTypeMap[tygoDollar[2].object.Name] = tygoDollar[2].object
		}
	case 5:
		tygoDollar = tygoS[tygopt-5 : tygopt+1]
		//line parser.y:76
		{
			eiota = 0
			tygoVAL.enum = &Enum{Name: tygoDollar[2].ident, Package: parserPkg, Values: make(map[string]int)}
		}
	case 6:
		tygoDollar = tygoS[tygopt-6 : tygopt+1]
		//line parser.y:81
		{
			tygoVAL.enum = tygoDollar[1].enum
			tygoVAL.enum.Values[tygoDollar[3].ident] = tygoDollar[5].integer
			eiota++
		}
	case 7:
		tygoDollar = tygoS[tygopt-6 : tygopt+1]
		//line parser.y:87
		{
			tygoVAL.enum = tygoDollar[1].enum
			tygoVAL.enum.Values[tygoDollar[3].ident] = eiota
			eiota++
		}
	case 8:
		tygoDollar = tygoS[tygopt-4 : tygopt+1]
		//line parser.y:93
		{
			tygoVAL.enum = tygoDollar[1].enum
			tygoVAL.enum.Values[tygoDollar[3].ident] = eiota
			eiota++
		}
	case 9:
		tygoDollar = tygoS[tygopt-5 : tygopt+1]
		//line parser.y:101
		{
			tygoVAL.object = &Object{Name: tygoDollar[2].ident, Package: parserPkg}
		}
	case 10:
		tygoDollar = tygoS[tygopt-5 : tygopt+1]
		//line parser.y:105
		{
			tygoVAL.object = tygoDollar[1].object
			tygoVAL.object.Fields = append(tygoVAL.object.Fields, &Field{Type: tygoDollar[4].spec, Name: tygoDollar[3].ident})
		}
	case 11:
		tygoDollar = tygoS[tygopt-6 : tygopt+1]
		//line parser.y:110
		{
			tygoVAL.object = tygoDollar[1].object
			tygoVAL.object.Fields = append(tygoVAL.object.Fields, &Field{Type: tygoDollar[5].spec, Name: tygoDollar[4].ident, Hide: true})
		}
	case 12:
		tygoDollar = tygoS[tygopt-4 : tygopt+1]
		//line parser.y:115
		{
			if tygoDollar[1].object.Parent != nil {
				log.Fatalf("[Tygo][Parser] Multiple inheritance is not allowed!")
			}
			tygoVAL.object = tygoDollar[1].object
			tygoVAL.object.Parent = tygoDollar[3].spec.(*InstanceType)
		}
	case 13:
		tygoDollar = tygoS[tygopt-3 : tygopt+1]
		//line parser.y:123
		{
			tygoVAL.object = tygoDollar[1].object
			tygoVAL.object.Methods = append(tygoVAL.object.Methods, tygoDollar[3].method)
		}
	case 15:
		tygoDollar = tygoS[tygopt-3 : tygopt+1]
		//line parser.y:131
		{
			tygoVAL.method = tygoDollar[1].method
			tygoVAL.method.Results = []Type{tygoDollar[2].spec}
		}
	case 16:
		tygoDollar = tygoS[tygopt-5 : tygopt+1]
		//line parser.y:136
		{
			tygoVAL.method = tygoDollar[1].method
			tygoVAL.method.Results = tygoDollar[3].specs
		}
	case 18:
		tygoDollar = tygoS[tygopt-3 : tygopt+1]
		//line parser.y:144
		{
			tygoVAL.method = tygoDollar[1].method
			tygoVAL.method.Params = []Type{tygoDollar[2].spec}
		}
	case 19:
		tygoDollar = tygoS[tygopt-3 : tygopt+1]
		//line parser.y:149
		{
			tygoVAL.method = tygoDollar[1].method
			tygoVAL.method.Params = tygoDollar[2].specs
		}
	case 20:
		tygoDollar = tygoS[tygopt-2 : tygopt+1]
		//line parser.y:156
		{
			tygoVAL.method = &Method{Name: tygoDollar[1].ident}
		}
	case 21:
		tygoDollar = tygoS[tygopt-3 : tygopt+1]
		//line parser.y:162
		{
			tygoVAL.specs = []Type{tygoDollar[1].spec, tygoDollar[3].spec}
		}
	case 22:
		tygoDollar = tygoS[tygopt-3 : tygopt+1]
		//line parser.y:166
		{
			tygoVAL.specs = append(tygoDollar[1].specs, tygoDollar[3].spec)
		}
	case 24:
		tygoDollar = tygoS[tygopt-3 : tygopt+1]
		//line parser.y:173
		{
			tygoVAL.spec = &ListType{E: tygoDollar[3].spec}
		}
	case 25:
		tygoDollar = tygoS[tygopt-5 : tygopt+1]
		//line parser.y:177
		{
			tygoVAL.spec = &DictType{K: tygoDollar[3].spec, V: tygoDollar[5].spec}
		}
	case 26:
		tygoDollar = tygoS[tygopt-4 : tygopt+1]
		//line parser.y:181
		{
			tygoVAL.spec = &VariantType{Ts: tygoDollar[3].specs}
		}
	case 27:
		tygoDollar = tygoS[tygopt-6 : tygopt+1]
		//line parser.y:185
		{
			tygoVAL.spec = &FixedPointType{Precision: uint(tygoDollar[3].integer), Floor: tygoDollar[5].integer}
		}
	case 28:
		tygoDollar = tygoS[tygopt-1 : tygopt+1]
		//line parser.y:191
		{
			if t, ok := parserTypeMap[tygoDollar[1].ident]; ok {
				switch i := t.(type) {
				case *Enum:
					tygoVAL.spec = &EnumType{Enum: i, Name: tygoDollar[1].ident}
				case *Object:
					tygoVAL.spec = &InstanceType{Object: i, Name: tygoDollar[1].ident}
				default:
					log.Fatalf("[Tygo][InstanceType] Unexpect type: %v", t)
				}
			} else if pkg, ok := parserTypePkg[tygoDollar[1].ident]; ok {
				tygoVAL.spec = &InstanceType{PkgName: pkg[0], PkgPath: pkg[1], Name: tygoDollar[1].ident}
			} else {
				tygoVAL.spec = SimpleType_FromString(tygoDollar[1].ident)
			}
		}
	case 29:
		tygoDollar = tygoS[tygopt-3 : tygopt+1]
		//line parser.y:208
		{
			tygoVAL.spec = &InstanceType{PkgName: tygoDollar[1].ident, PkgPath: parserImports[tygoDollar[1].ident], Name: tygoDollar[3].ident}
		}
	case 30:
		tygoDollar = tygoS[tygopt-2 : tygopt+1]
		//line parser.y:212
		{
			if t, ok := parserTypeMap[tygoDollar[2].ident]; ok {
				tygoVAL.spec = &InstanceType{Object: t.(*Object), IsPtr: true, Name: tygoDollar[2].ident}
			} else if pkg, ok := parserTypePkg[tygoDollar[2].ident]; ok {
				tygoVAL.spec = &InstanceType{IsPtr: true, PkgName: pkg[0], PkgPath: pkg[1], Name: tygoDollar[2].ident}
			} else {
				tygoVAL.spec = &InstanceType{IsPtr: true, Name: tygoDollar[2].ident}
			}
		}
	case 31:
		tygoDollar = tygoS[tygopt-4 : tygopt+1]
		//line parser.y:222
		{
			tygoVAL.spec = &InstanceType{IsPtr: true, PkgName: tygoDollar[2].ident, PkgPath: parserImports[tygoDollar[2].ident], Name: tygoDollar[4].ident}
		}
	}
	goto tygostack /* stack new state and value */
}

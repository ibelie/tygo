//line decorator.y:10
package tygo

import __yyfmt__ "fmt"

//line decorator.y:11
import (
	"bytes"
	"log"
	"unicode/utf8"
)

//line decorator.y:21
type decoratorSymType struct {
	yys        int
	id         string
	decorator  *Decorator
	decorators []*Decorator
}

const ID = 57346

var decoratorToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"'@'",
	"','",
	"'('",
	"')'",
	"ID",
}
var decoratorStatenames = [...]string{}

const decoratorEofCode = 1
const decoratorErrCode = 2
const decoratorInitialStackSize = 16

//line decorator.y:66

// The parser expects the lexer to return 0 on EOF.  Give it a name for clarity.
const eof = 0

// The parser uses the type <prefix>Lex as a lexer. It must provide
// the methods Lex(*<prefix>SymType) int and Error(string).
type decoratorLex struct {
	line []byte
	peek rune
}

// The parser calls this method to get each new token. This
// implementation returns symbols and ID.
func (x *decoratorLex) Lex(yylval *decoratorSymType) int {
	for {
		c := x.next()
		switch c {
		case eof:
			return eof
		case '@', ',', '(', ')':
			return int(c)
		case ' ', '\t', '\n', '\r':
		default:
			return x.id(c, yylval)
		}
	}
}

// Lex a id.
func (x *decoratorLex) id(c rune, yylval *decoratorSymType) int {
	add := func(b *bytes.Buffer, c rune) {
		if _, err := b.WriteRune(c); err != nil {
			log.Fatalf("[Tygo][Decorator] WriteRune: %s", err)
		}
	}
	var b bytes.Buffer
	add(&b, c)
L:
	for {
		c = x.next()
		switch c {
		case '@', ',', '(', ')', ' ', '\t', '\n', '\r':
			break L
		default:
			add(&b, c)
		}
	}
	if c != eof {
		x.peek = c
	}
	yylval.id = b.String()
	return ID
}

// Return the next rune for the lexer.
func (x *decoratorLex) next() rune {
	if x.peek != eof {
		r := x.peek
		x.peek = eof
		return r
	}
	if len(x.line) == 0 {
		return eof
	}
	c, size := utf8.DecodeRune(x.line)
	x.line = x.line[size:]
	if c == utf8.RuneError && size == 1 {
		log.Print("[Tygo][Decorator] Invalid utf8")
		return x.next()
	}
	return c
}

// The parser calls this method on a parse error.
func (x *decoratorLex) Error(s string) {
	log.Printf("[Tygo][Decorator] Parse error: %s", s)
}

//line yacctab:1
var decoratorExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const decoratorNprod = 7
const decoratorPrivate = 57344

var decoratorTokenNames []string
var decoratorStates []string

const decoratorLast = 13

var decoratorAct = [...]int{

	4, 5, 7, 11, 6, 10, 3, 2, 9, 8,
	1, 0, 12,
}
var decoratorPact = [...]int{

	3, 2, -7, -7, -1000, -4, -1000, -7, -2, -1000,
	-1000, -7, -1000,
}
var decoratorPgo = [...]int{

	0, 0, 10, 9,
}
var decoratorR1 = [...]int{

	0, 2, 2, 1, 1, 3, 3,
}
var decoratorR2 = [...]int{

	0, 2, 3, 1, 4, 1, 3,
}
var decoratorChk = [...]int{

	-1000, -2, 4, 4, -1, 8, -1, 6, -3, -1,
	7, 5, -1,
}
var decoratorDef = [...]int{

	0, -2, 0, 0, 1, 3, 2, 0, 0, 5,
	4, 0, 6,
}
var decoratorTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	6, 7, 3, 3, 5, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 4,
}
var decoratorTok2 = [...]int{

	2, 3, 8,
}
var decoratorTok3 = [...]int{
	0,
}

var decoratorErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	decoratorDebug        = 0
	decoratorErrorVerbose = false
)

type decoratorLexer interface {
	Lex(lval *decoratorSymType) int
	Error(s string)
}

type decoratorParser interface {
	Parse(decoratorLexer) int
	Lookahead() int
}

type decoratorParserImpl struct {
	lval  decoratorSymType
	stack [decoratorInitialStackSize]decoratorSymType
	char  int
}

func (p *decoratorParserImpl) Lookahead() int {
	return p.char
}

func decoratorNewParser() decoratorParser {
	return &decoratorParserImpl{}
}

const decoratorFlag = -1000

func decoratorTokname(c int) string {
	if c >= 1 && c-1 < len(decoratorToknames) {
		if decoratorToknames[c-1] != "" {
			return decoratorToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func decoratorStatname(s int) string {
	if s >= 0 && s < len(decoratorStatenames) {
		if decoratorStatenames[s] != "" {
			return decoratorStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func decoratorErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !decoratorErrorVerbose {
		return "syntax error"
	}

	for _, e := range decoratorErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + decoratorTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := decoratorPact[state]
	for tok := TOKSTART; tok-1 < len(decoratorToknames); tok++ {
		if n := base + tok; n >= 0 && n < decoratorLast && decoratorChk[decoratorAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if decoratorDef[state] == -2 {
		i := 0
		for decoratorExca[i] != -1 || decoratorExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; decoratorExca[i] >= 0; i += 2 {
			tok := decoratorExca[i]
			if tok < TOKSTART || decoratorExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if decoratorExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += decoratorTokname(tok)
	}
	return res
}

func decoratorlex1(lex decoratorLexer, lval *decoratorSymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = decoratorTok1[0]
		goto out
	}
	if char < len(decoratorTok1) {
		token = decoratorTok1[char]
		goto out
	}
	if char >= decoratorPrivate {
		if char < decoratorPrivate+len(decoratorTok2) {
			token = decoratorTok2[char-decoratorPrivate]
			goto out
		}
	}
	for i := 0; i < len(decoratorTok3); i += 2 {
		token = decoratorTok3[i+0]
		if token == char {
			token = decoratorTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = decoratorTok2[1] /* unknown char */
	}
	if decoratorDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", decoratorTokname(token), uint(char))
	}
	return char, token
}

func decoratorParse(decoratorlex decoratorLexer) int {
	return decoratorNewParser().Parse(decoratorlex)
}

func (decoratorrcvr *decoratorParserImpl) Parse(decoratorlex decoratorLexer) int {
	var decoratorn int
	var decoratorVAL decoratorSymType
	var decoratorDollar []decoratorSymType
	_ = decoratorDollar // silence set and not used
	decoratorS := decoratorrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	decoratorstate := 0
	decoratorrcvr.char = -1
	decoratortoken := -1 // decoratorrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		decoratorstate = -1
		decoratorrcvr.char = -1
		decoratortoken = -1
	}()
	decoratorp := -1
	goto decoratorstack

ret0:
	return 0

ret1:
	return 1

decoratorstack:
	/* put a state and value onto the stack */
	if decoratorDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", decoratorTokname(decoratortoken), decoratorStatname(decoratorstate))
	}

	decoratorp++
	if decoratorp >= len(decoratorS) {
		nyys := make([]decoratorSymType, len(decoratorS)*2)
		copy(nyys, decoratorS)
		decoratorS = nyys
	}
	decoratorS[decoratorp] = decoratorVAL
	decoratorS[decoratorp].yys = decoratorstate

decoratornewstate:
	decoratorn = decoratorPact[decoratorstate]
	if decoratorn <= decoratorFlag {
		goto decoratordefault /* simple state */
	}
	if decoratorrcvr.char < 0 {
		decoratorrcvr.char, decoratortoken = decoratorlex1(decoratorlex, &decoratorrcvr.lval)
	}
	decoratorn += decoratortoken
	if decoratorn < 0 || decoratorn >= decoratorLast {
		goto decoratordefault
	}
	decoratorn = decoratorAct[decoratorn]
	if decoratorChk[decoratorn] == decoratortoken { /* valid shift */
		decoratorrcvr.char = -1
		decoratortoken = -1
		decoratorVAL = decoratorrcvr.lval
		decoratorstate = decoratorn
		if Errflag > 0 {
			Errflag--
		}
		goto decoratorstack
	}

decoratordefault:
	/* default state action */
	decoratorn = decoratorDef[decoratorstate]
	if decoratorn == -2 {
		if decoratorrcvr.char < 0 {
			decoratorrcvr.char, decoratortoken = decoratorlex1(decoratorlex, &decoratorrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if decoratorExca[xi+0] == -1 && decoratorExca[xi+1] == decoratorstate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			decoratorn = decoratorExca[xi+0]
			if decoratorn < 0 || decoratorn == decoratortoken {
				break
			}
		}
		decoratorn = decoratorExca[xi+1]
		if decoratorn < 0 {
			goto ret0
		}
	}
	if decoratorn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			decoratorlex.Error(decoratorErrorMessage(decoratorstate, decoratortoken))
			Nerrs++
			if decoratorDebug >= 1 {
				__yyfmt__.Printf("%s", decoratorStatname(decoratorstate))
				__yyfmt__.Printf(" saw %s\n", decoratorTokname(decoratortoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for decoratorp >= 0 {
				decoratorn = decoratorPact[decoratorS[decoratorp].yys] + decoratorErrCode
				if decoratorn >= 0 && decoratorn < decoratorLast {
					decoratorstate = decoratorAct[decoratorn] /* simulate a shift of "error" */
					if decoratorChk[decoratorstate] == decoratorErrCode {
						goto decoratorstack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if decoratorDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", decoratorS[decoratorp].yys)
				}
				decoratorp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if decoratorDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", decoratorTokname(decoratortoken))
			}
			if decoratortoken == decoratorEofCode {
				goto ret1
			}
			decoratorrcvr.char = -1
			decoratortoken = -1
			goto decoratornewstate /* try again in the same state */
		}
	}

	/* reduction by production decoratorn */
	if decoratorDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", decoratorn, decoratorStatname(decoratorstate))
	}

	decoratornt := decoratorn
	decoratorpt := decoratorp
	_ = decoratorpt // guard against "declared and not used"

	decoratorp -= decoratorR2[decoratorn]
	// decoratorp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if decoratorp+1 >= len(decoratorS) {
		nyys := make([]decoratorSymType, len(decoratorS)*2)
		copy(nyys, decoratorS)
		decoratorS = nyys
	}
	decoratorVAL = decoratorS[decoratorp+1]

	/* consult goto table to find next state */
	decoratorn = decoratorR1[decoratorn]
	decoratorg := decoratorPgo[decoratorn]
	decoratorj := decoratorg + decoratorS[decoratorp].yys + 1

	if decoratorj >= decoratorLast {
		decoratorstate = decoratorAct[decoratorg]
	} else {
		decoratorstate = decoratorAct[decoratorj]
		if decoratorChk[decoratorstate] != -decoratorn {
			decoratorstate = decoratorAct[decoratorg]
		}
	}
	// dummy call; replaced with literal code
	switch decoratornt {

	case 1:
		decoratorDollar = decoratorS[decoratorpt-2 : decoratorpt+1]
		//line decorator.y:37
		{
			decoratorVAL.decorators = []*Decorator{decoratorDollar[2].decorator}
		}
	case 2:
		decoratorDollar = decoratorS[decoratorpt-3 : decoratorpt+1]
		//line decorator.y:41
		{
			decoratorVAL.decorators = append(decoratorDollar[1].decorators, decoratorDollar[3].decorator)
		}
	case 3:
		decoratorDollar = decoratorS[decoratorpt-1 : decoratorpt+1]
		//line decorator.y:47
		{
			decoratorVAL.decorator = &Decorator{Name: decoratorDollar[1].id}
		}
	case 4:
		decoratorDollar = decoratorS[decoratorpt-4 : decoratorpt+1]
		//line decorator.y:51
		{
			decoratorVAL.decorator = &Decorator{Name: decoratorDollar[1].id, Params: decoratorDollar[3].decorators}
		}
	case 5:
		decoratorDollar = decoratorS[decoratorpt-1 : decoratorpt+1]
		//line decorator.y:57
		{
			decoratorVAL.decorators = []*Decorator{decoratorDollar[1].decorator}
		}
	case 6:
		decoratorDollar = decoratorS[decoratorpt-3 : decoratorpt+1]
		//line decorator.y:61
		{
			decoratorVAL.decorators = append(decoratorDollar[1].decorators, decoratorDollar[3].decorator)
		}
	}
	goto decoratorstack /* stack new state and value */
}

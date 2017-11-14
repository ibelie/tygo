// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

// This is a parser of tygo.
// To build it:
// go tool yacc -o "parser.y.go" -p "tygo" parser.y

%{

package tygo

import (
	"bytes"
	"log"
	"strconv"
	"unicode/utf8"
)

%}

%union {
	ident   string
	keyword string
	integer int
	spec    Type
	specs   []Type
	method  *Method
	object  *Object
	enum    *Enum
}

%type	<enum>    enum
%type	<object>  object
%type	<method>  method method1 method2
%type	<specs>   specs
%type	<spec>    spec spec1

%token	'*' '=' '.' ',' '(' ')' '[' ']' '{' '}' '<' '>' '\t' '\n'
%token	<keyword> TYPE ENUM OBJECT MAP FIXEDPOINT VARIANT IOTA
%token	<ident>   IDENT
%token	<integer> INTEGER

%%

top:
	enum '}' newline
	{
		parserTypes = append(parserTypes, $1)
		parserTypeMap[$1.Name] = $1
	}
|	top enum '}' newline
	{
		parserTypes = append(parserTypes, $2)
		parserTypeMap[$2.Name] = $2
	}
|	object '}' newline
	{
		if $1.Parent == nil {
			$1.Parent = &InstanceType{PkgName: "tygo", PkgPath: TYGO_PATH, Name: "Tygo"}
		}
		parserTypes = append(parserTypes, $1)
		parserTypeMap[$1.Name] = $1
	}
|	top object '}' newline
	{
		if $2.Parent == nil {
			$2.Parent = &InstanceType{PkgName: "tygo", PkgPath: TYGO_PATH, Name: "Tygo"}
		}
		parserTypes = append(parserTypes, $2)
		parserTypeMap[$2.Name] = $2
	}

enum:
	TYPE IDENT ENUM '{' newline
	{
		eiota = 0
		$$ = &Enum{Name: $2, Package: parserPkg, Values: make(map[string]int)}
	}
|	enum '\t' IDENT '=' INTEGER newline
	{
		$$ = $1
		$$.Values[$3] = $5
		eiota++
	}
|	enum '\t' IDENT '=' IOTA newline
	{
		$$ = $1
		$$.Values[$3] = eiota
		eiota++
	}
|	enum '\t' IDENT newline
	{
		$$ = $1
		$$.Values[$3] = eiota
		eiota++
	}

object:
	TYPE IDENT OBJECT '{' newline
	{
		$$ = &Object{Name: $2, Package: parserPkg}
	}
|	object '\t' IDENT spec newline
	{
		$$ = $1
		$$.Fields = append($$.Fields, &Field{Type: $4, Name: $3})
	}
|	object '\t' '.' IDENT spec newline
	{
		$$ = $1
		$$.Fields = append($$.Fields, &Field{Type: $5, Name: $4, Hide: true})
	}
|	object '\t' spec1 newline
	{
		if $1.Parent != nil {
			log.Fatalf("[Tygo][Parser] Multiple inheritance is not allowed!")
		}
		$$ = $1
		$$.Parent = $3.(*InstanceType)
	}
|	object '\t' method
	{
		$$ = $1
		$$.Methods = append($$.Methods, $3)
	}

method:
	method1 newline
|	method1 spec newline
	{
		$$ = $1
		$$.Results = []Type{$2}
	}
|	method1 '(' specs ')' newline
	{
		$$ = $1
		$$.Results = $3
	}

method1:
	method2 ')'
|	method2 spec ')'
	{
		$$ = $1
		$$.Params = []Type{$2}
	}
|	method2 specs ')'
	{
		$$ = $1
		$$.Params = $2
	}

method2:
	IDENT '('
	{
		$$ = &Method{Name: $1}
	}

specs:
	spec ',' spec
	{
		$$ = []Type{$1, $3}
	}
|	specs ',' spec
	{
		$$ = append($1, $3)
	}

spec:
	spec1
|	'[' ']' spec
	{
		$$ = &ListType{E: $3}
	}
|	MAP '[' spec ']' spec
	{
		$$ = &DictType{K: $3, V: $5}
	}
|	VARIANT '<' specs '>'
	{
		$$ = &VariantType{Ts: $3}
	}
|	FIXEDPOINT '<' INTEGER ',' INTEGER '>'
	{
		$$ = &FixedPointType{Precision: uint($3), Floor: $5}
	}

spec1:
	IDENT
	{
		if t, ok := parserTypeMap[$1]; ok {
			switch i := t.(type) {
			case *Enum:
				$$ = &EnumType{Enum: i, Name: $1}
			case *Object:
				$$ = &InstanceType{Object: i, Name: $1}
			default:
				log.Fatalf("[Tygo][InstanceType] Unexpect type: %v", t)
			}
		} else if pkg, ok := parserTypePkg[$1]; ok {
			$$ = &InstanceType{PkgName: pkg[0], PkgPath: pkg[1], Name: $1}
		} else {
			$$ = SimpleType_FromString($1)
		}
	}
|	IDENT '.' IDENT
	{
		$$ = &InstanceType{PkgName: $1, PkgPath: parserImports[$1], Name: $3}
	}
|	'*' IDENT
	{
		if t, ok := parserTypeMap[$2]; ok {
			$$ = &InstanceType{Object: t.(*Object), IsPtr: true, Name: $2}
		} else if pkg, ok := parserTypePkg[$2]; ok {
			$$ = &InstanceType{IsPtr: true, PkgName: pkg[0], PkgPath: pkg[1], Name: $2}
		} else {
			$$ = &InstanceType{IsPtr: true, Name: $2}
		}
	}
|	'*' IDENT '.' IDENT
	{
		$$ = &InstanceType{IsPtr: true, PkgName: $2, PkgPath: parserImports[$2], Name: $4}
	}

newline:
	'\n'
|	newline '\n'


%%

var eiota int

var (
	parserPkg     string
	parserTypes   []Type
	parserTypeMap map[string]Type
	parserImports map[string]string
	parserTypePkg map[string][2]string
)

func Parse(code string, pkg string, imports map[string]string, typePkg map[string][2]string) ([]Type) {
	parserPkg     = pkg
	parserTypes   = nil
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
	L: for {
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
	L: for {
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

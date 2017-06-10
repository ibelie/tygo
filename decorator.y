// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

// This is a parser of tygo decorator.
// To build it:
// go tool yacc -o "decorator.y.go" -p "decorator" decorator.y

%{

package tygo

import (
	"bytes"
	"log"
	"unicode/utf8"
)

%}

%union {
	id            string
	decorator    *Decorator
	decorators []*Decorator
}

%type	<decorator>  dec
%type	<decorators> top decs

%token	'@' ',' '(' ')'
%token	<id> ID

%%

top:
	'@' dec
	{
		$$ = []*Decorator{$2}
	}
|	top '@' dec
	{
		$$ = append($1, $3)
	}

dec:
	ID
	{
		$$ = &Decorator{Name: $1}
	}
|	ID '(' decs ')'
	{
		$$ = &Decorator{Name: $1, Params: $3}
	}

decs:
	dec
	{
		$$ = []*Decorator{$1}
	}
|	decs ',' dec
	{
		$$ = append($1, $3)
	}


%%

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
	L: for {
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

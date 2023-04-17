package tealex

import (
	"strings"
	"unicode/utf8"
)

type Lexer struct {
	p int // previous index
	i int // current index

	l  int // current line index
	lb int // current line begin index in source
	li int // current index in the current line

	ts  []Token // tokens
	tsi int     // tokens index

	diag []Error // errors

	Source []byte
}

func (l *Lexer) fail(msg string) {
	l.diag = append(l.diag, Error{
		l:  l.l,
		li: l.li,

		msg: msg,
	})
}

func (l *Lexer) inc(n int) {
	l.i += n
	l.li += n
}

func (l *Lexer) emit(t TokenType) {
	l.ts = append(l.ts, Token{
		l: l.l,
		b: l.p - l.lb,
		e: l.i - l.lb,

		v: string(l.Source[l.p:l.i]),
		t: t,
	})

	l.p = l.i

	switch t {
	case TokenEol:
		l.l++
		l.li = 0
		l.lb = l.i
	}
}

func isTerminating(c rune) bool {
	switch c {
	case '\r':
	case '\n':
	case ' ':
	case '\t':
	case ';':
	default:
		return false
	}

	return true
}

func (z *Lexer) readValue() {
	p, n := utf8.DecodeRune(z.Source[z.i:])
	if p == '"' {
		z.inc(n)
		for {
			if z.i == len(z.Source) {
				z.fail("incomplete string")
				return
			}

			c, n := utf8.DecodeRune(z.Source[z.i:])
			if c == '"' && p != '\\' {
				z.inc(n)

				s := string(z.Source[z.p+1 : z.i-1])
				v := "\"" + strings.ReplaceAll(s, "\\\"", "\"") + "\""

				z.ts = append(z.ts, Token{
					l: z.l,
					b: z.p - z.lb,
					e: z.i - z.lb,

					v: v,
					t: TokenValue,
				})

				z.p = z.i
				return
			}

			z.inc(n)

			p = c
		}
	} else {
		for {
			if z.i == len(z.Source) {
				z.emit(TokenValue)
				return
			}

			c, n := utf8.DecodeRune(z.Source[z.i:])
			if isTerminating(c) {
				z.emit(TokenValue)
				return
			}

			z.inc(n)
		}
	}
}

func (z *Lexer) skipWhitespace() {
	for {
		if z.i == len(z.Source) {
			return
		}

		c, n := utf8.DecodeRune(z.Source[z.i:])
		if c != ' ' && c != '\t' {
			return
		}

		z.inc(n)
		z.p = z.i
	}
}

func (z *Lexer) readComment() {
	for {
		l := z.i - z.p
		if z.i == len(z.Source) {
			if l < 2 {
				z.fail("incomplete comment")
				return
			}

			z.ts = append(z.ts, Token{
				l: z.l,
				b: z.p - z.lb,
				e: z.i - z.lb,

				v: string(z.Source[z.p+2 : z.i]),
				t: TokenComment,
			})

			z.p = z.i
			return
		}

		c, n := utf8.DecodeRune(z.Source[z.i:])
		if l < 2 {
			if c != '/' {
				z.fail("incomplete comment")
				return
			}
		} else {
			if c == '\r' || c == '\n' {
				z.ts = append(z.ts, Token{
					l: z.l,
					b: z.p - z.lb,
					e: z.i - z.lb,

					v: string(z.Source[z.p+2 : z.i]),
					t: TokenComment,
				})

				z.p = z.i
				return
			}
		}

		z.inc(n)
	}
}

func (z *Lexer) readSeparator() {
	for {
		if z.i == len(z.Source) {
			z.emit(TokenEol)
			return
		}

		c, n := utf8.DecodeRune(z.Source[z.i:])
		switch c {
		case '\r':
			z.inc(n)
			if z.i < len(z.Source) {
				c2, n2 := utf8.DecodeRune(z.Source[z.i:])
				if c2 == '\n' {
					z.inc(n2)
				}
			}
			z.emit(TokenEol)
			return
		case '\n':
			z.inc(n)
			z.emit(TokenEol)
			return
		case ';':
			z.inc(n)
			z.emit(TokenSemicolon)
			return
		}

		z.inc(n)
	}
}

func (z *Lexer) readTokens() {
	z.skipWhitespace()

	if z.i < len(z.Source) {
		c, n := utf8.DecodeRune(z.Source[z.i:])

		var nc rune
		if z.i+n < len(z.Source) {
			nc, _ = utf8.DecodeRune(z.Source[z.i+n:])
		}

		if c == '/' && nc == '/' {
			z.readComment()
		} else if c == '\n' || c == '\r' || c == ';' {
			z.readSeparator()
		} else {
			z.readValue()
		}
	}
}

func (z *Lexer) read() {
	if z.i == len(z.Source) {
		return
	}

	z.readTokens()
}

func (z *Lexer) Return() bool {
	if z.tsi == 0 {
		return false
	}

	z.tsi--

	return true
}

func (z *Lexer) Scan() bool {
	if len(z.ts) > z.tsi {
		z.tsi++
	}

	if len(z.ts) == z.tsi {
		z.read()
	}

	return len(z.ts) > z.tsi
}

func (z *Lexer) Curr() Token {
	if len(z.ts) == z.tsi {
		return Token{}
	}

	return z.ts[z.tsi]
}

func (z *Lexer) Errors() []Error {
	return z.diag
}

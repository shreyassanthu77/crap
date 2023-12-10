package parser

import (
	"fmt"
	"os"
)

type Loc struct {
	Line int
	Col  int
	Pos  int
}

type ParseError struct {
	Msg string
	Loc Loc
}

func (e *ParseError) Throw() {
	fmt.Printf("%d:%d: ParseError: %s\n", e.Loc.Line, e.Loc.Col, e.Msg)
	os.Exit(1)
}

type Parser struct {
	input  string
	length int
	pos    int
	line   int
	col    int
}

type Rule struct {
	Selector Selector
	Decls    []Declaration
}

type Stylesheet []Rule

func Parse(input string) Stylesheet {
	return (&Parser{
		input:  input,
		length: len(input),
		pos:    0,
		line:   1,
		col:    0,
	}).parse()
}

func (p *Parser) loc() Loc {
	return Loc{
		Line: p.line,
		Col:  p.col,
		Pos:  p.pos,
	}
}

func (p *Parser) hasNext(n ...int) bool {
	if len(n) > 0 {
		return p.pos+n[0] < p.length
	}
	return p.pos < p.length
}

func (p *Parser) peek(by ...int) byte {
	n := 0
	if len(by) > 0 {
		n = by[0]
	}
	if !p.hasNext(n) {
		return 0
	}
	return p.input[p.pos+n]
}

func (p *Parser) next() byte {
	if !p.hasNext() {
		p.err("(next) unexpected end of input")
	}
	c := p.input[p.pos]
	if c == '\n' || (c == '\r' && p.peek() != '\n') {
		p.line++
		p.col = 1
	} else {
		p.col++
	}
	p.pos++
	return c
}

func (p *Parser) skipWhitespace() {
	for p.hasNext() && isWhitespace(p.peek()) {
		p.next()
	}
}

func (p *Parser) expect(chs ...byte) *ParseError {
	var err *ParseError
	for _, ch := range chs {
		if !p.hasNext() {
			return p.err("unexpected end of input")
		}
		if p.peek() != ch {
			err = p.err(fmt.Sprintf("expected '%c' got '%c'", ch, p.peek()))
		} else {
			err = nil
			p.next()
			break
		}
	}
	return err
}

func (p *Parser) err(msg string) *ParseError {
	return &ParseError{
		Msg: msg,
		Loc: p.loc(),
	}
}

func (p *Parser) parseString() (string, *ParseError) {
	start := p.pos
	ch := p.peek()
	if ch != '"' && ch != '\'' {
		return "", p.err("expected string")
	}
	p.next()
	quote := ch
	for {
		ch = p.peek()
		if ch == '\\' && p.peek(1) == quote {
			p.next()
		}
		if ch == quote {
			p.next()
			break
		}
		p.next()
	}
	return p.input[start:p.pos], nil
}

func (p *Parser) parseNumber() (string, *ParseError) {
	start := p.pos
	ch := p.peek()
	if !isDigit(ch) && ch != '.' {
		return "", p.err("expected number")
	}
	p.next()
	for {
		ch = p.peek()
		if !isDigit(ch) && ch != '.' {
			break
		}
		p.next()
	}
	return p.input[start:p.pos], nil
}

func (p *Parser) parse() Stylesheet {
	css := Stylesheet{}
	for {
		selector, err := p.parseSelector()
		if err != nil {
			break
		}
		decls, err := p.parseBlock()
		if err != nil {
			err.Throw()
		}

		css = append(css, Rule{
			Selector: selector,
			Decls:    decls,
		})
	}
	return css
}

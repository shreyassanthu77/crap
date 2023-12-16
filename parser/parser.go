package parser

import (
	. "github.con/shreyascodes-tech/sss-lang/ast"
	"github.con/shreyascodes-tech/sss-lang/lexer"
)

type Parser struct {
	input  string
	lexer  *lexer.Lexer
	tokens []lexer.Token
	pos    int
}

func New(input string) *Parser {
	return &Parser{
		input: input,
		lexer: lexer.New(input),
	}
}

func (p *Parser) Parse() Stylesheet {
	stylesheet := Stylesheet{}

	i := 0
	for {
		i++
		if i > 100 {
			panic("infinite loop")
		}
		token := p.next()
		switch token.Type {
		case lexer.EOF:
			return stylesheet
		case lexer.COMMENT:
			continue
		case lexer.AT:
			rule, ok := p.parseAtRule(token)
			if !ok {
				continue
			}
			stylesheet.AtRules = append(stylesheet.AtRules, rule)
		}
	}
}

func (p *Parser) parseAtRule(token lexer.Token) (AtRule, bool) {
	rule := AtRule{}
	rule.Loc.Start = p.span(token.Loc.Start)

	id, ok := p.parseIdentifier()
	if !ok {
		return AtRule{}, false
	}
	rule.Name = id

	value, ok := p.parseValue()
	if !ok {
		return AtRule{}, false
	}
	rule.Param = value

	semi, ok := p.eat(lexer.SEMICOLON)
	if ok {
		rule.Loc.End = p.span(semi.Loc.End)
	}

	return rule, true
}

func (p *Parser) parseIdentifier() (Identifier, bool) {
	if token, ok := p.eat(lexer.IDENT); ok {
		return Identifier{
			Loc:   p.loc(token),
			Value: token.Value,
		}, true
	}
	return Identifier{}, false
}

func (p *Parser) parseValue() (Value, bool) {
	token := p.next()
	switch token.Type {
	case lexer.STRING:
		return Value{
			Loc:  p.loc(token),
			Type: STRING,
			Data: token.Value,
		}, true
	case lexer.NUMBER:
		return Value{
			Loc:  p.loc(token),
			Type: NUMBER,
			Data: float64(0),
		}, true
	case lexer.HEX:
		return Value{
			Loc:  p.loc(token),
			Type: HEX,
			Data: token.Value,
		}, true
	case lexer.IDENT:
		id := Value{
			Loc:  p.loc(token),
			Type: IDENT,
			Data: token.Value,
		}
		if p.peek().Type == lexer.LPAREN {
			return p.parseFunctionCall(id)
		}
		return id, true
		// case lexer.LSQUIRLY:
		// 	return p.parseExpression()
	}
	return Value{}, false
}

func (p *Parser) parseFunctionCall(id Value) (Value, bool) {
	p.eat(lexer.LPAREN)

	res := Value{
		Loc:  id.Loc,
		Type: FUNCTION_CALL,
	}
	fn := FunctionCall{
		Loc: id.Loc,
		Name: Identifier{
			Loc:   id.Loc,
			Value: id.Data.(string),
		},
		Args: []Value{},
	}

	for {
		next := p.peek()
		if next.Type == lexer.RPAREN {
			p.next()
			fn.Loc.End = p.span(next.Loc.End)
			break
		}
		arg, ok := p.parseValue()
		if !ok {
			return Value{}, false
		}
		fn.Args = append(fn.Args, arg)

		if _, ok := p.eat(lexer.COMMA); !ok && p.peek().Type != lexer.RPAREN {
			return Value{}, false
		}
	}

	res.Data = fn
	res.Loc.End = fn.Loc.End
	return res, true
}

func (p *Parser) span(token lexer.Span) Span {
	return Span{
		Pos:  token.Pos,
		Line: token.Line,
		Col:  token.Col,
	}
}

func (p *Parser) loc(token lexer.Token) Location {
	return Location{
		Start: p.span(token.Loc.Start),
		End:   p.span(token.Loc.End),
	}
}

func (p *Parser) eat(t lexer.TokenType) (lexer.Token, bool) {
	token := p.next()
	if token.Type == t {
		return token, true
	} else {
		p.rewind(1)
		return token, false
	}
}

func (p *Parser) peek() lexer.Token {
	if p.pos >= len(p.tokens) {
		next := p.lexer.Next()
		if next.Type == lexer.EOF {
			return next
		} else {
			p.tokens = append(p.tokens, next)
		}
	}
	tok := p.tokens[p.pos]
	return tok
}

func (p *Parser) next() lexer.Token {
	next := p.peek()
	if next.Type != lexer.EOF {
		p.pos++
	}
	return next
}

func (p *Parser) rewind(by int) {
	p.pos -= by
}

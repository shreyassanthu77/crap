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

	rules := p.parserRules()
	stylesheet.Rules = append(stylesheet.Rules, rules...)

	if len(rules) > 0 {
		stylesheet.Loc.Start = rules[0].Loc.Start
		stylesheet.Loc.End = rules[len(rules)-1].Loc.End
	}

	return stylesheet
}

func (p *Parser) parserRules() []RuleNode {
	rules := []RuleNode{}

loop:
	for {
		next := p.next()
		switch next.Type {
		case lexer.COMMENT:
			continue
		case lexer.EOF:
			break loop
		case lexer.AT:
			rule, ok := p.parseAtRule(next)
			if !ok {
				break loop
			}
			rules = append(rules, rule)
		case lexer.COLON:
			p.rewind(1)
			rule, ok := p.parseRule()
			if !ok {
				break loop
			}
			rules = append(rules, rule)
		case lexer.IDENT:
			if p.peek().Type == lexer.COLON {
				rule, ok := p.parseDeclaration(next)
				if !ok {
					break loop
				}
				rules = append(rules, rule)
				continue
			}
			p.rewind(1)
			rule, ok := p.parseRule()
			if !ok {
				break loop
			}
			rules = append(rules, rule)
		default:
			p.rewind(1)
			break loop
		}
	}
	return rules
}

func (p *Parser) parseDeclaration(token lexer.Token) (RuleNode, bool) {
	rule := RuleNode{}
	declaration := Declaration{}
	rule.Loc.Start = p.span(token.Loc.Start)

	declaration.Name = Identifier{
		Loc:   p.loc(token),
		Value: token.Value,
	}

	if _, ok := p.eat(lexer.COLON); !ok {
		return rule, false
	}

	value, ok := p.parseValue()
	if !ok {
		return rule, false
	}
	declaration.Value = value

	if _, ok := p.eat(lexer.SEMICOLON); ok {
		rule.Loc.End = p.span(token.Loc.End)
	}

	rule.Rule = declaration
	rule.Type = DECL_RULE
	return rule, true
}

func (p *Parser) parseAtRule(token lexer.Token) (RuleNode, bool) {
	rule := RuleNode{}
	atRule := AtRule{}
	rule.Loc.Start = p.span(token.Loc.Start)

	id, ok := p.parseIdentifier()
	if !ok {
		return rule, false
	}
	atRule.Name = id

	value, ok := p.parseValue()
	if !ok {
		return rule, false
	}

	atRule.Params = []Value{value}

	if lSquirly, ok := p.eat(lexer.LSQUIRLY); ok {
		children, loc := p.parseBlock(lSquirly)
		atRule.Rules = children
		rule.Loc.End = loc.End
	}

	if _, ok := p.eat(lexer.SEMICOLON); ok {
		rule.Loc.End = p.span(token.Loc.End)
	}

	rule.Rule = atRule
	rule.Type = AT_RULE
	return rule, true
}

func (p *Parser) parseRule() (RuleNode, bool) {
	rule := RuleNode{}
	styleRule := StyleRule{}

	selectors, ok := p.parseSelectors()
	if !ok {
		return rule, false
	}
	styleRule.Selectors = selectors

	if lSquirly, ok := p.eat(lexer.LSQUIRLY); ok {
		children, loc := p.parseBlock(lSquirly)
		styleRule.Rules = children
		rule.Loc.End = loc.End
	}

	rule.Rule = styleRule
	rule.Type = STYLE_RULE
	return rule, true
}

func (p *Parser) parseSelectors() ([]Selector, bool) {
	selectors := []Selector{}

	for {
		next := p.peek()
		switch next.Type {
		case lexer.COMMA:
			p.next()
			continue
		case lexer.LSQUIRLY:
			return selectors, true
		case lexer.EOF:
			return selectors, true
		}

		selector, ok := p.parseSelector()
		if !ok {
			return selectors, false
		}
		selectors = append(selectors, selector)
	}
}

func (p *Parser) parseSelector() (Selector, bool) {
	selector := Selector{}

loop:
	for {
		next := p.peek()
		switch next.Type {
		case lexer.COMMA:
			return selector, true
		case lexer.LSQUIRLY:
			return selector, true
		case lexer.EOF:
			return selector, true
		default:
			if isCombinator(next) {
				if selector.Next != nil {
					return selector, false
				}
				selector.Combinator = next.Value
				p.next()
				next, ok := p.parseSelector()
				if ok {
					selector.Next = &next
					break loop
				}
			}

			id, ok := p.parseSelectorPart()
			if !ok {
				break loop
			}
			selector.Parts = append(selector.Parts, id)
		}
	}

	return selector, true
}

func (p *Parser) parseSelectorPart() (SelectorPart, bool) {
	next := p.peek()
	part := SelectorPart{}
	part.Loc.Start = p.span(next.Loc.Start)

	switch next.Type {
	case lexer.IDENT:
		part.Type = TYPE_SELECTOR
	case lexer.DOT:
		p.next()
		part.Type = CLASS_SELECTOR
	case lexer.OCTOTHORPE:
		p.next()
		part.Type = ID_SELECTOR
	case lexer.LBRACKET:
		part.Type = ATTRIBUTE_SELECTOR
	case lexer.COLON:
		p.next()
		part.Type = PSEUDO_SELECTOR
	}

	id, ok := p.parseIdentifier()

	if !ok {
		return part, false
	}

	part.Value = Value{
		Loc:  id.Loc,
		Type: IDENT,
		Data: id,
	}

	part.Loc.End = id.Loc.End
	return part, true
}

func (p *Parser) parseBlock(lSquirly lexer.Token) ([]RuleNode, Location) {
	rules := p.parserRules()

	rSquirly, ok := p.eat(lexer.RSQUIRLY)
	if !ok {
		return nil, Location{}
	}

	return rules, Location{
		Start: p.span(lSquirly.Loc.Start),
		End:   p.span(rSquirly.Loc.End),
	}
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
		p.rewind(1)
		id, _ := p.parseIdentifier()
		if p.peek().Type == lexer.LPAREN {
			return p.parseFunctionCall(id)
		}
		return Value{
			Loc:  id.Loc,
			Type: IDENT,
			Data: id,
		}, true
		// case lexer.LSQUIRLY:
		// 	return p.parseExpression()
	}
	return Value{}, false
}

func (p *Parser) parseFunctionCall(id Identifier) (Value, bool) {
	p.eat(lexer.LPAREN)

	res := Value{
		Loc:  id.Loc,
		Type: FUNCTION_CALL,
	}
	fn := FunctionCall{
		Loc:  id.Loc,
		Name: id,
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

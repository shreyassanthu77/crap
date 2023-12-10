package parser

type Declaration struct {
	Property string
	Value    string
}

func (p *Parser) parseProperty() (string, *ParseError) {
	start := p.pos
	ch := p.peek()
	if !isAlpha(ch) && ch != '-' {
		return "", p.err("expected property name")
	}
	p.next()
	ch = p.peek()
	for {
		if !isAlphaNumeric(ch) && ch != '-' {
			break
		}
		p.next()
		ch = p.peek()
	}
	return p.input[start:p.pos], nil
}

func (p *Parser) parseValue() (string, *ParseError) {
	str, err := p.parseString()
	if err == nil {
		return str, nil
	}
	num, err := p.parseNumber()
	if err == nil {
		return num, nil
	}
	return "", p.err("expected string or number")
}

func (p *Parser) parseDeclaration() (Declaration, *ParseError) {
	rule := Declaration{}
	p.skipWhitespace()
	property, err := p.parseProperty()
	if err != nil {
		return rule, err
	}
	rule.Property = property
	p.skipWhitespace()
	if err := p.expect(':'); err != nil {
		return rule, err
	}
	p.skipWhitespace()
	value, err := p.parseValue()
	if err != nil {
		return rule, err
	}
	rule.Value = value
	p.skipWhitespace()
	if err := p.expect(';'); err != nil {
		return rule, err
	}
	return rule, nil
}

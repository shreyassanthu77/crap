package parser

type Selector struct {
	Tag       string
	Attribute map[string]string
}

func (p *Parser) parseTag() (string, *ParseError) {
	p.skipWhitespace()

	start := p.pos
	ch := p.peek()
	if !isAlpha(ch) && ch != '*' && ch != '.' && ch != '#' && ch != ':' {
		return "", p.err("expected tag name")
	}
	p.next()
	for {
		ch = p.peek()
		if !isAlpha(ch) {
			break
		}
		p.next()
	}
	return p.input[start:p.pos], nil
}

func (p *Parser) parseAttributeKey() (string, *ParseError) {
	start := p.pos
	ch := p.peek()
	if !isAlpha(ch) {
		return "", p.err("expected key name")
	}
	ch = p.next()
	for {
		ch = p.peek()
		if !isAlpha(ch) && ch != '-' {
			break
		}
		p.next()
	}
	return p.input[start:p.pos], nil
}

func (p *Parser) parseAttributeValue() (string, *ParseError) {
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

func (p *Parser) parseAttributes() (map[string]string, *ParseError) {
	attributes := make(map[string]string)
	for {
		p.skipWhitespace()
		if err := p.expect('['); err != nil {
			break
		}
		p.skipWhitespace()
		key, err := p.parseAttributeKey()
		if err != nil {
			return nil, err
		}
		p.skipWhitespace()
		if err := p.expect('='); err != nil {
			return nil, err
		}
		p.skipWhitespace()
		value, err := p.parseAttributeValue()
		if err != nil {
			return nil, err
		}
		p.skipWhitespace()
		if err := p.expect(']'); err != nil {
			return nil, err
		}
		attributes[key] = value
	}
	return attributes, nil
}

func (p *Parser) parseSelector() (Selector, *ParseError) {
	selector := Selector{}
	tag, err := p.parseTag()
	if err != nil {
		return selector, err
	}
	selector.Tag = tag
	attribute, err := p.parseAttributes()
	selector.Attribute = attribute
	return selector, err
}

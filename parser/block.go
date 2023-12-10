package parser

func (p *Parser) parseBlock() ([]Declaration, *ParseError) {
	ruleset := []Declaration{}
	p.skipWhitespace()
	if err := p.expect('{'); err != nil {
		return nil, err
	}
	for {
		p.skipWhitespace()
		rule, err := p.parseDeclaration()
		if err != nil {
			break
		}
		ruleset = append(ruleset, rule)
		p.skipWhitespace()
	}
	if err := p.expect('}'); err != nil {
		return nil, err
	}
	return ruleset, nil
}

package parser

import (
	. "github.com/shreyassanthu77/cisp/ast"
	"github.com/shreyassanthu77/cisp/lexer"
)

func (p *Parser) parseSelector(ident Identifier) (Selector, error) {
	next, err := p.peek()
	if err != nil {
		return Selector{}, err
	}

	if next.Typ == lexer.TOK_LBRACKET {
		attr, err := p.parseAttributes()
		if err != nil {
			return Selector{}, err
		}

		return Selector{
			Identifier: ident,
			Atrributes: attr,
		}, nil
	}

	next, err = p.peek()
	if err != nil {
		return Selector{}, err
	}

	if next.Typ == lexer.TOK_IDENTIFIER {
		panic("Complex selectors not implemented")
	}

	return Selector{
		Identifier: ident,
		Atrributes: nil,
	}, nil
}

func (p *Parser) parseAttributes() ([]Attreibute, error) {
	attrs := []Attreibute{}

	for {
		next, err := p.peek()
		if err != nil {
			return nil, err
		}

		if next.Typ != lexer.TOK_LBRACKET {
			break
		}

		p.next() // Consume '['
		id, err := p.expect(lexer.TOK_IDENTIFIER)
		if err != nil {
			return nil, err
		}

		next, err = p.peek()
		if err != nil {
			return nil, err
		}

		attr := Attreibute{
			Name:    Identifier{Name: id.Value},
			Default: NilValue{},
		}

		if next.Typ == lexer.TOK_EQUAL {
			p.next() // Consume '='
			val, err := p.parseLiteralVal()
			if err != nil {
				return nil, err
			}

			attr.Default = val
		}
		_, err = p.expect(lexer.TOK_RBRACKET)
		if err != nil {
			return nil, err
		}
	}

	return attrs, nil
}

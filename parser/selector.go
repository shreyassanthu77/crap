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
		Atrributes: map[Identifier]Value{},
	}, nil
}

func (p *Parser) parseAttributes() (map[Identifier]Value, error) {
	attrs := map[Identifier]Value{}

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

		if next.Typ == lexer.TOK_EQUAL {
			p.next() // Consume '='
			val, err := p.parseValue()
			if err != nil {
				return nil, err
			}

			attrs[Identifier{Name: id.Value}] = val
		} else {
			attrs[Identifier{Name: id.Value}] = NilValue{}
		}

		_, err = p.expect(lexer.TOK_RBRACKET)
		if err != nil {
			return nil, err
		}
	}

	return attrs, nil
}

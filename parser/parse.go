package parser

import (
	. "github.com/shreyassanthu77/cisp/ast"
	"github.com/shreyassanthu77/cisp/lexer"
)

func (p *Parser) parseFunctionCall(name string) (FunctionCall, error) {
	_, err := p.expect(lexer.TOK_LPAREN)
	if err != nil {
		return FunctionCall{}, err
	}

	params := []Value{}
	next, err := p.peek()
	if err != nil {
		return FunctionCall{}, err
	}

	if next.Typ == lexer.TOK_RPAREN {
		p.next()
		return FunctionCall{
			Name:       name,
			Parameters: params,
		}, nil
	}

	for {
		param, err := p.parseValue()
		if err != nil {
			return FunctionCall{}, err
		}
		params = append(params, param)

		next, err = p.peek()
		if err != nil {
			return FunctionCall{}, err
		}

		if next.Typ != lexer.TOK_COMMA {
			break
		}
		p.next()
		next, err = p.peek()
	}

	_, err = p.expect(lexer.TOK_RPAREN)
	if err != nil {
		return FunctionCall{}, err
	}

	return FunctionCall{
		Name:       name,
		Parameters: params,
	}, nil
}

func (p *Parser) parseDeclarationStmt(id Identifier) (Declaration, error) {
	_, err := p.expect(lexer.TOK_COLON)
	if err != nil {
		return Declaration{
			Property: id,
		}, err
	}

	values := []Value{}
	for {
		next, err := p.peek()
		if err != nil {
			return Declaration{}, err
		}

		if next.Typ == lexer.TOK_SEMICOLON {
			break
		}

		value, err := p.parseValue()
		if err != nil {
			return Declaration{}, err
		}

		values = append(values, value)
	}

	_, err = p.expect(lexer.TOK_SEMICOLON)
	if err != nil {
		return Declaration{}, err
	}

	return Declaration{
		Property: id,
		Values:   values,
	}, nil
}

func (p *Parser) parseStatement() (Statement, error) {
	next, err := p.peek()
	if err != nil {
		return nil, err
	}

	if next.Typ == lexer.TOK_AT {
		return p.parseAtRule()
	}

	id, err := p.expect(lexer.TOK_IDENTIFIER)
	if err != nil {
		return nil, err
	}

	next, err = p.peek()
	if err != nil {
		return nil, err
	}

	if next.Typ == lexer.TOK_COLON {
		return p.parseDeclarationStmt(Identifier{Name: id.Value})
	}

	return p.parseNestedRule(Identifier{Name: id.Value})
}

func (p *Parser) parseDeclarationBlock() ([]Statement, error) {
	_, err := p.expect(lexer.TOK_LSQUIRLY)
	if err != nil {
		return nil, err
	}

	stmts := []Statement{}
	var lastErr error
	for {
		next, err := p.peek()
		if err != nil {
			return nil, err
		}

		if next.Typ == lexer.TOK_RSQUIRLY {
			break
		}

		stmt, err := p.parseStatement()
		if err != nil {
			lastErr = err
			break
		}

		stmts = append(stmts, stmt)
	}

	if lastErr != nil {
		return nil, lastErr
	}

	_, err = p.expect(lexer.TOK_RSQUIRLY)
	if err != nil {
		return nil, err
	}

	return stmts, nil
}

func (p *Parser) parseNestedRule(id Identifier) (Rule, error) {
	selector, err := p.parseSelector(id)
	if err != nil {
		return Rule{}, err
	}

	decls, err := p.parseDeclarationBlock()
	if err != nil {
		return Rule{}, err
	}

	return Rule{
		Selector:   selector,
		Statements: decls,
	}, nil
}

func (p *Parser) parseRule() (Rule, error) {
	id, err := p.expect(lexer.TOK_IDENTIFIER)
	if err != nil {
		return Rule{}, err
	}

	return p.parseNestedRule(Identifier{Name: id.Value})
}

func (p *Parser) parseAtRule() (AtRule, error) {
	p.next() // Consume '@'
	name, err := p.expect(lexer.TOK_IDENTIFIER)
	if err != nil {
		return AtRule{}, err
	}

	params := []Value{}
	for {
		next, err := p.peek()
		if err != nil {
			return AtRule{}, err
		}

		if next.Typ == lexer.TOK_LSQUIRLY || next.Typ == lexer.TOK_SEMICOLON {
			break
		}

		param, err := p.parseValue()
		if err != nil {
			return AtRule{}, err
		}

		params = append(params, param)
	}

	next, err := p.peek()
	if err != nil {
		return AtRule{}, err
	}

	if next.Typ == lexer.TOK_LSQUIRLY {
		decls, err := p.parseDeclarationBlock()
		if err != nil {
			return AtRule{}, err
		}

		return AtRule{
			Name:       Identifier{Name: name.Value},
			Params:     params,
			Statements: decls,
		}, nil
	}

	_, err = p.expect(lexer.TOK_SEMICOLON)
	if err != nil {
		return AtRule{}, err
	}

	return AtRule{
		Name:   Identifier{Name: name.Value},
		Params: params,
	}, nil
}

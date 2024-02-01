package parser

import (
	"fmt"
	"strconv"

	. "github.com/shreyassanthu77/cisp/ast"
	"github.com/shreyassanthu77/cisp/lexer"
)

type Parser struct {
	lex    *lexer.Lexer
	tok    lexer.Token
	hasTok bool
}

func New(lex *lexer.Lexer) *Parser {
	return &Parser{
		lex:    lex,
		hasTok: false,
	}
}

func (p *Parser) peek() (lexer.Token, error) {
	if p.hasTok {
		return p.tok, nil
	}

	tok, err := p.lex.Next()
	if err != nil {
		return lexer.Token{}, err
	}

	p.tok = tok
	p.hasTok = true
	return tok, nil
}

func (p *Parser) next() (lexer.Token, error) {
	tok, err := p.peek()
	if err != nil {
		return lexer.Token{}, err
	}
	p.hasTok = false
	return tok, nil
}

func (p *Parser) expect(typ string) (lexer.Token, error) {
	tok, err := p.next()
	if err != nil {
		return lexer.Token{}, err
	}
	if tok.Typ != typ {
		return lexer.Token{}, fmt.Errorf("%d:%d Expected %s but got %s", tok.Line, tok.Col, typ, tok.Typ)
	}
	return tok, nil
}

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

func (p *Parser) parseValue() (Value, error) {
	tok, err := p.next()
	if err != nil {
		return nil, err
	}

	switch tok.Typ {
	case lexer.TOK_IDENTIFIER:
		next, err := p.peek()
		if err != nil {
			return nil, err
		}
		if next.Typ == lexer.TOK_LPAREN {
			return p.parseFunctionCall(tok.Value)
		}
		return Identifier{Name: tok.Value}, nil
	case lexer.TOK_STRING:
		return String{Value: tok.Value}, nil
	case lexer.TOK_INT:
		f, err := strconv.ParseInt(tok.Value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse number: %s", err)
		}
		return Int{Value: f}, nil
	case lexer.TOK_FLOAT:
		f, err := strconv.ParseFloat(tok.Value, 64)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse number: %s", err)
		}
		return Float{Value: f}, nil
	case lexer.TOK_TRUE:
		return Boolean{Value: true}, nil
	case lexer.TOK_FALSE:
		return Boolean{Value: false}, nil
	}

	fmt.Println(tok)
	panic("Unimplemented")
}

func (p *Parser) parseDeclarationStmt(id Identifier) (Declaration, error) {
	_, err := p.expect(lexer.TOK_COLON)
	if err != nil {
		return Declaration{
			Property: id,
		}, err
	}

	value, err := p.parseValue()
	if err != nil {
		return Declaration{}, err
	}

	_, err = p.expect(lexer.TOK_SEMICOLON)
	if err != nil {
		return Declaration{}, err
	}

	return Declaration{
		Property: id,
		Val:      value,
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

func (p *Parser) Parse() (Program, error) {
	rules := []IRule{}
	var lastErr error
outer:
	for {
		next, err := p.peek()

		// Skip semicolons
		for {
			if err != nil {
				return Program{}, err
			}

			if next.Typ == lexer.EOF {
				break outer
			}

			if next.Typ == lexer.TOK_SEMICOLON {
				p.next()
				next, err = p.peek()
				continue
			}
			break
		}

		if next.Typ == lexer.TOK_AT {
			atRule, err := p.parseAtRule()
			if err != nil {
				lastErr = err
				break
			}
			rules = append(rules, atRule)
		} else {
			rule, err := p.parseRule()
			if err != nil {
				lastErr = err
				break
			}
			rules = append(rules, rule)
		}
	}

	if lastErr != nil {
		return Program{}, lastErr
	}
	return Program{
		Rules: rules,
	}, nil
}

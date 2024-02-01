package parser

import (
	"fmt"
	"strconv"

	. "github.com/shreyassanthu77/cisp/ast"
	"github.com/shreyassanthu77/cisp/lexer"
)

func (p *Parser) parseLiteralVal() (Value, error) {
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
	case lexer.TOK_EMPTY:
		return NilValue{}, nil
	}

	fmt.Println(tok)
	panic("Unimplemented")
}

func (p *Parser) parseValue() (Value, error) {
	return p.parseLiteralVal()
}

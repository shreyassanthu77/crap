package main

import (
	"fmt"
	"strconv"

	"github.com/shreyassanthu77/cisp/lexer"
)

/*
	Syntax

	id = [a-zA-Z_][a-zA-Z0-9_]*
	string = ('"' + .* + '"') | ("'" + .* + "'")
	number = (# + [a-fA-F0-9]{3,6,8})
		| ([0-9]* + (. + [0-9]+)?)
	boolean = true | false
	value = string | number | boolean | function_call
	unary_operator = ! | ~
	operator = + | - | * | / | % | ^ | = | == | != | > | >= | < | <= | && | ||
	expression = value | unary_operator + expression
		| function_call
		| (expression + operator + expression)
		| '(' + expression + ')'
	function_call = id + '(' + function_parameters + ')'
	function_parameters = expression + (, + expression)* + (,)*

	selector = (. | #)? + id + ([id=value])*
		+ (selector_delimiter + selector)*
	selector_delimiter = (> | + | ~)

	at_rule = @ + id + expression + (delcaration_block | ;)
	rule = (selector + declaration_block) | at_rule
	delcaration_block = { + (declaration | rule)* + }
	declaration = id + : + value;

	program = (rule | at_rule)*
*/

type Parser struct {
	lex    *lexer.Lexer
	tok    lexer.Token
	hasTok bool
}

func NewParser(lex *lexer.Lexer) *Parser {
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

type Value interface {
	isValue()
}

type Identifier struct {
	Name string
}

func (i Identifier) isValue() {}

func (i Identifier) String() string {
	return fmt.Sprintf("Id(%s)", i.Name)
}

type String struct {
	Value string
}

func (s String) isValue() {}

func (s String) String() string {
	return fmt.Sprintf("'%s'", s.Value)
}

type Number struct {
	Value float64
}

func (n Number) isValue() {}

func (n Number) String() string {
	return fmt.Sprintf("%f", n.Value)
}

type Boolean struct {
	Value bool
}

func (b Boolean) isValue() {}

func (b Boolean) String() string {
	return fmt.Sprintf("Bool(%t)", b.Value)
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
		return Identifier{tok.Value}, nil
	case lexer.TOK_STRING:
		return String{tok.Value}, nil
	case lexer.TOK_NUMBER:
		f, err := strconv.ParseFloat(tok.Value, 64)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse number: %s", err)
		}
		return Number{f}, nil
	case lexer.TOK_TRUE:
		return Boolean{true}, nil
	case lexer.TOK_FALSE:
		return Boolean{false}, nil
	}

	panic("Unimplemented")
}

type FunctionCall struct {
	Name       string
	Parameters []Value
}

func (f FunctionCall) isValue() {}

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

type Declaration struct {
	Property Identifier
	Val      Value
}

func (p *Parser) parseDeclaration() (Declaration, error) {
	id, err := p.expect(lexer.TOK_IDENTIFIER)
	if err != nil {
		return Declaration{}, err
	}

	_, err = p.expect(lexer.TOK_COLON)
	if err != nil {
		return Declaration{}, err
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
		Property: Identifier{id.Value},
		Val:      value,
	}, nil
}

func (p *Parser) parseDeclarationBlock() ([]Declaration, error) {
	_, err := p.expect(lexer.TOK_LSQUIRLY)
	if err != nil {
		return nil, err
	}

	decls := []Declaration{}
	var lastErr error
	for {
		next, err := p.peek()
		if err != nil {
			return nil, err
		}

		if next.Typ == lexer.TOK_RSQUIRLY {
			break
		}

		decl, err := p.parseDeclaration()
		if err != nil {
			lastErr = err
			break
		}

		decls = append(decls, decl)
	}

	if lastErr != nil {
		return nil, lastErr
	}

	_, err = p.expect(lexer.TOK_RSQUIRLY)
	if err != nil {
		return nil, err
	}

	return decls, nil
}

func (p *Parser) Parse() error {
	panic("Unimplemented")
}

func main() {
	input := `{
		a: 10;
		b: "Hello World";
		c: #fff;
		d: true;
		e: false;
		f: foo();
		g: foo(10, "hello", true, #fff, foo());
}`
	lex := lexer.New(input)
	parser := NewParser(lex)

	decl, err := parser.parseDeclarationBlock()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%+v\n", decl)
}

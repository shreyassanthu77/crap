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
	number = ([0-9]* + (. + [0-9]+)?)
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

	selector = (. | #)? + id + ([id=value])* + (, + selector)*

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

type NilValue struct{}

func (n NilValue) isValue() {}

func (n NilValue) String() string {
	return "Nil"
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

type Int struct {
	Value int64
}

func (n Int) isValue() {}

func (n Int) String() string {
	return fmt.Sprintf("%d", n.Value)
}

type Float struct {
	Value float64
}

func (n Float) isValue() {}

func (n Float) String() string {
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
	case lexer.TOK_INT:
		f, err := strconv.ParseInt(tok.Value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse number: %s", err)
		}
		return Int{f}, nil
	case lexer.TOK_FLOAT:
		f, err := strconv.ParseFloat(tok.Value, 64)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse number: %s", err)
		}
		return Float{f}, nil
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

func (f FunctionCall) String() string {
	return fmt.Sprintf("Call(%s, %v)", f.Name, f.Parameters)
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

type Statement interface {
	isStatement()
}

type Declaration struct {
	Property Identifier
	Val      Value
}

func (d Declaration) isStatement() {}

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
	id, err := p.expect(lexer.TOK_IDENTIFIER)
	if err != nil {
		return nil, err
	}

	next, err := p.peek()
	if err != nil {
		return nil, err
	}

	if next.Typ == lexer.TOK_COLON {
		return p.parseDeclarationStmt(Identifier{id.Value})
	}

	return p.parseNestedRule(Identifier{id.Value})
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

type Selector struct {
	Identifier Identifier
	Atrributes map[Identifier]Value
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

			attrs[Identifier{id.Value}] = val
		} else {
			attrs[Identifier{id.Value}] = NilValue{}
		}

		_, err = p.expect(lexer.TOK_RBRACKET)
		if err != nil {
			return nil, err
		}
	}

	return attrs, nil
}

type IRule interface {
	isRule()
}

type Rule struct {
	Selector   Selector
	Statements []Statement
}

func (r Rule) isRule() {}

func (r Rule) isStatement() {}

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

	return p.parseNestedRule(Identifier{id.Value})
}

type AtRule struct {
	Name       Identifier
	Params     []Value
	Statements []Statement
}

func (r AtRule) isRule() {}

func (r AtRule) isStatement() {}

func (p *Parser) parseAtRule() (AtRule, error) {
	p.next() // Consume '@'

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
			Name:       Identifier{next.Value},
			Params:     params,
			Statements: decls,
		}, nil
	}

	_, err = p.expect(lexer.TOK_SEMICOLON)
	if err != nil {
		return AtRule{}, err
	}

	return AtRule{
		Name:   Identifier{next.Value},
		Params: params,
	}, nil
}

func (p *Parser) Parse() ([]IRule, error) {
	rules := []IRule{}
	var lastErr error
	for {
		next, err := p.peek()
		if err != nil {
			return nil, err
		}

		if next.Typ == lexer.EOF {
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
		return nil, lastErr
	}

	return rules, nil
}

func main() {
	input := `@hello world {
			hi {
				foo: bar;
		}
	}
	@hello world2;

	hello {
		foo: bar;
	}`
	lex := lexer.New(input)
	parser := NewParser(lex)

	decl, err := parser.Parse()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%+v\n", decl)
}

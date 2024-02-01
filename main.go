package main

import (
	"fmt"
	"github.com/shreyassanthu77/cisp/lexer"
	"github.com/shreyassanthu77/cisp/parser"
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
	par := parser.New(lex)

	decl, err := par.Parse()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%+v\n", decl)
}

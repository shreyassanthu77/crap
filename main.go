package main

import (
	"fmt"
)

/*
	Syntax

	id = [a-zA-Z_][a-zA-Z0-9_]*
	value = id | string | number | boolean | expression
	string = '"' + .* + '"'
	number = (# + [a-fA-F0-9]{3,6,8})
		| ([0-9]* + (. + [0-9]+)?)
	boolean = true | false
	unary_operator = ! | ~
	operator = + | - | * | / | % | ^ | = | == | != | > | >= | < | <= | && | ||
	expression =
		| unary_operator + expression
		| function_call
		| (expression + operator + expression)
		| '(' + expression + ')'
	function_call = id + '(' + function_parameters + ')'
	function_parameters = expression + (, + expression)* + (,)*

	selector = (. | #)? + id + ([id=value])*
		+ (selector_delimiter + selector)*
	selector_delimiter = (> | + | ~)

	at_rule = @ + id + value* + (delcaration_block | ;)
	rule = (selector + declaration_block) | at_rule
	delcaration_block = { + (declaration | rule)* + }
	declaration = id + : + value + important? + ;

	stylesheet = (rule | at_rule)*
	important = "!important"
*/

const (
	// Token types
	TOK_EOF        = "EOF"
	TOK_IDENTIFIER = "IDENTIFIER"
	TOK_STRING     = "STRING"
	TOK_NUMBER     = "NUMBER"
	TOK_BOOLEAN    = "BOOLEAN"

	// Unary operators
	TOK_BANG  = "!"
	TOK_TILDE = "~"

	// Binary operators
	TOK_PLUS               = "+"
	TOK_MINUS              = "-"
	TOK_ASTERISK           = "*"
	TOK_SLASH              = "/"
	TOK_PERCENT            = "%"
	TOK_CARET              = "^"
	TOK_EQUAL              = "="
	TOK_DOUBLE_EQUAL       = "=="
	TOK_NOT_EQUAL          = "!="
	TOK_GREATER_THAN       = ">"
	TOK_GREATER_THAN_EQUAL = ">="
	TOK_LESS_THAN          = "<"
	TOK_LESS_THAN_EQUAL    = "<="
	TOK_AND                = "&&"
	TOK_OR                 = "||"

	// Expressions
	TOK_LPAREN = "("
	TOK_RPAREN = ")"
	TOK_COMMA  = ","

	// Selectors
	TOK_DOT  = "."
	TOK_HASH = "#"

	// rules
	TOK_AT        = "@"
	TOK_LSQUIRLY  = "{"
	TOK_RSQUIRLY  = "}"
	TOK_COLON     = ":"
	TOK_SEMICOLON = ";"
	TOK_IMPORTANT = "!important"
)

func main() {
	fmt.Println("Hello World!")
}

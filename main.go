package main

import (
	"fmt"
)

/*
	Syntax

	id = [a-zA-Z_][a-zA-Z0-9_]*
	value = id | string | number | boolean | expression
	string = ('"' + .* + '"') | ("'" + .* + "'")
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

// Lexer
const (
	EOF            = "EOF"
	TOK_IDENTIFIER = "IDENTIFIER"
	TOK_STRING     = "STRING"
	TOK_NUMBER     = "NUMBER"
	TOK_BOOLEAN    = "BOOLEAN"

	// Unary operators
	TOK_BANG  = "BANG"
	TOK_TILDE = "TILDE"

	// Binary operators
	TOK_PLUS               = "PLUS"
	TOK_MINUS              = "MINUS"
	TOK_ASTERISK           = "ASTERISK"
	TOK_SLASH              = "SLASH"
	TOK_PERCENT            = "PERCENT"
	TOK_CARET              = "CARET"
	TOK_EQUAL              = "EQUAL"
	TOK_DOUBLE_EQUAL       = "DOUBLE_EQUAL"
	TOK_NOT_EQUAL          = "NOT_EQUAL"
	TOK_GREATER_THAN       = "GREATER_THAN"
	TOK_GREATER_THAN_EQUAL = "GREATER_THAN_EQUAL"
	TOK_LESS_THAN          = "LESS_THAN"
	TOK_LESS_THAN_EQUAL    = "LESS_THAN_EQUAL"
	TOK_AND                = "AND"
	TOK_OR                 = "OR"

	// Expressions
	TOK_LPAREN = "LPAREN"
	TOK_RPAREN = "RPAREN"
	TOK_COMMA  = "COMMA"

	// Selectors
	TOK_DOT  = "DOT"
	TOK_HASH = "HASH"

	// rules
	TOK_AT        = "AT"
	TOK_LSQUIRLY  = "LSQUIRLY"
	TOK_RSQUIRLY  = "RSQUIRLY"
	TOK_COLON     = "COLON"
	TOK_SEMICOLON = "SEMICOLON"
	TOK_IMPORTANT = "IMPORTANT"
)

type Token struct {
	typ   string
	value string
}

func (t Token) String() string {
	return fmt.Sprintf("TOK<%s>(%s)", t.typ, t.value)
}

type Lexer struct {
	input string
	pos   int
	line  int
	col   int
	done  bool
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input: input,
		pos:   0,
		line:  1,
		col:   1,
		done:  false,
	}
}

func (l *Lexer) peek() string {
	if l.done {
		return EOF
	}

	if l.pos >= len(l.input) {
		l.done = true
		return EOF
	}

	return string(l.input[l.pos])
}

func (l *Lexer) next() string {
	ch := l.peek()
	if ch == EOF {
		return EOF
	}

	l.pos++
	l.col++

	nextCh := l.peek()
	if ch == "\n" || (ch == "\r" && nextCh != "\n") {
		l.line++
		l.col = 1
	}

	return ch
}

func (l *Lexer) skipWhitespace() {
	if l.done {
		return
	}
	for {
		ch := l.peek()
		if ch == EOF {
			break
		}

		if ch != " " && ch != "\t" && ch != "\n" && ch != "\r" {
			break
		}

		l.next()
	}
}

func (l *Lexer) readString(quote string) (Token, error) {
	if l.done {
		return Token{}, fmt.Errorf("Unexpected EOF")
	}

	start := l.pos
	for {
		ch := l.next()
		if ch == EOF {
			return Token{}, fmt.Errorf("Unexpected EOF")
		}

		if ch == quote {
			break
		}
	}

	str := l.input[start : l.pos-1]
	return Token{typ: TOK_STRING, value: str}, nil
}

func (l *Lexer) Next() (Token, error) {
	l.skipWhitespace()
	ch := l.next()

	if ch == EOF {
		return Token{typ: EOF}, nil
	}
	nextCh := l.peek()

	switch ch {
	case "!":
		if nextCh == "=" {
			l.next()
			return Token{typ: TOK_NOT_EQUAL, value: ch + nextCh}, nil
		}
		return Token{typ: TOK_BANG, value: ch}, nil
	case "~":
		return Token{typ: TOK_TILDE, value: ch}, nil
	case "+":
		return Token{typ: TOK_PLUS, value: ch}, nil
	case "-":
		return Token{typ: TOK_MINUS, value: ch}, nil
	case "*":
		return Token{typ: TOK_ASTERISK, value: ch}, nil
	case "/":
		return Token{typ: TOK_SLASH, value: ch}, nil
	case "%":
		return Token{typ: TOK_PERCENT, value: ch}, nil
	case "^":
		return Token{typ: TOK_CARET, value: ch}, nil
	case "=":
		if nextCh == "=" {
			l.next()
			return Token{typ: TOK_DOUBLE_EQUAL, value: ch + nextCh}, nil
		}
		return Token{typ: TOK_EQUAL, value: ch}, nil
	case ">":
		if nextCh == "=" {
			l.next()
			return Token{typ: TOK_GREATER_THAN_EQUAL, value: ch + nextCh}, nil
		}
		return Token{typ: TOK_GREATER_THAN, value: ch}, nil
	case "<":
		if nextCh == "=" {
			l.next()
			return Token{typ: TOK_LESS_THAN_EQUAL, value: ch + nextCh}, nil
		}
		return Token{typ: TOK_LESS_THAN, value: ch}, nil
	case "&":
		if nextCh == "&" {
			l.next()
			return Token{typ: TOK_AND, value: ch + nextCh}, nil
		}
		return Token{}, fmt.Errorf("Unexpected character: %s", ch)
	case "|":
		if nextCh == "|" {
			l.next()
			return Token{typ: TOK_OR, value: ch + nextCh}, nil
		}
		return Token{}, fmt.Errorf("Unexpected character: %s", ch)
	case "(":
		return Token{typ: TOK_LPAREN, value: ch}, nil
	case ")":
		return Token{typ: TOK_RPAREN, value: ch}, nil
	case ",":
		return Token{typ: TOK_COMMA, value: ch}, nil
	case ".":
		return Token{typ: TOK_DOT, value: ch}, nil
	case "#":
		return Token{typ: TOK_HASH, value: ch}, nil
	case "@":
		return Token{typ: TOK_AT, value: ch}, nil
	case "{":
		return Token{typ: TOK_LSQUIRLY, value: ch}, nil
	case "}":
		return Token{typ: TOK_RSQUIRLY, value: ch}, nil
	case ":":
		return Token{typ: TOK_COLON, value: ch}, nil
	case ";":
		return Token{typ: TOK_SEMICOLON, value: ch}, nil
	case "\"":
		return l.readString(ch)
	case "'":
		return l.readString(ch)
	}

	panic("Not implemented")
}

func main() {
	test := `
		"hello" 'world'
		! ~ + - * / % ^ = == != > >= < <= && ||
		( ) , . #
		@ { } : ;
	`
	lexer := NewLexer(test)
	for !lexer.done {
		tok, err := lexer.Next()
		if err != nil {
			panic(err)
		}
		fmt.Println(tok)
	}
}

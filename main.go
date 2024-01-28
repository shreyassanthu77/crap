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
	TOK_TRUE       = "TRUE"
	TOK_FALSE      = "FALSE"

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
	pos   int
	line  int
	col   int
}

func (t Token) String() string {
	return fmt.Sprintf("{%d:%d} TOK<%s>(%s)", t.line, t.col, t.typ, t.value)
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

func (l *Lexer) tok(typ string, value string) Token {
	return Token{
		typ:   typ,
		value: value,
		pos:   l.pos,
		line:  l.line,
		col:   l.col - len(value) - 1,
	}
}

func (l *Lexer) error(format string, args ...interface{}) error {
	pre := fmt.Sprintf("{%d:%d} ", l.line, l.col)
	return fmt.Errorf(pre+format, args...)
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
	if ch == "\n" || ch == "\r" {
		l.line++
		l.col = 1
		if ch == "\r" && nextCh == "\n" {
			l.pos++
		}
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
		return Token{}, l.error("Unexpected EOF")
	}

	start := l.pos
	for {
		ch := l.next()
		if ch == EOF {
			return Token{}, l.error("Unexpected EOF")
		}

		if ch == quote {
			break
		}
	}

	str := l.input[start : l.pos-1]
	return l.tok(TOK_STRING, str), nil
}

func isValidHex(ch string) bool {
	return (ch >= "a" && ch <= "f") || (ch >= "A" && ch <= "F") || (ch >= "0" && ch <= "9")
}

func (l *Lexer) readHex() (Token, error) {
	if l.done {
		return Token{}, l.error("Unexpected EOF")
	}

	start := l.pos - 1 // -1 because we already read the `#`
	for {
		ch := l.next()
		if ch == EOF {
			return Token{}, l.error("Unexpected EOF")
		}

		if !isValidHex(ch) {
			break
		}
	}

	length := l.pos - start
	if length != 4 && length != 7 && length != 9 {
		l.pos = start // reset pos
		return Token{}, l.error("Invalid hex value")
	}

	str := l.input[start : l.pos-1]
	return l.tok(TOK_NUMBER, str), nil
}

func isValidIdentifierStart(ch string) bool {
	return ch == "_" || (ch >= "a" && ch <= "z") || (ch >= "A" && ch <= "Z")
}

func isValidIdentifier(ch string) bool {
	return isValidIdentifierStart(ch) || (ch >= "0" && ch <= "9") || ch == "-" || ch == "." || ch == "#"
}

func (l *Lexer) readIdentifier() (Token, error) {
	if l.done {
		return Token{}, l.error("Unexpected EOF")
	}

	start := l.pos - 1 // -1 because we already read the first char
	for {
		ch := l.peek()
		if ch == EOF {
			break
		}

		if !isValidIdentifier(ch) {
			break
		}

		l.next()
	}

	id := l.input[start:l.pos]

	switch id {
	case "true":
		return l.tok(TOK_TRUE, id), nil
	case "false":
		return l.tok(TOK_FALSE, id), nil
	}

	return l.tok(TOK_IDENTIFIER, id), nil
}

func isValidNumber(ch string) bool {
	return (ch >= "0" && ch <= "9") || (ch >= "a" && ch <= "f") || (ch >= "A" && ch <= "F")
}

func (l *Lexer) readNumber(ch string) (Token, error) {
	if l.done {
		return Token{}, l.error("Unexpected EOF")
	}

	start := l.pos - 1 // -1 because we already read the first char
	deci := ch == "."
	for {
		ch := l.peek()
		if ch == EOF {
			break
		}

		if ch == "." {
			if deci {
				l.error("Unexpected character: `.` after `.` in number is that a typo?")
			}
			deci = true
			l.next()
			continue
		}

		if !isValidNumber(ch) {
			break
		}

		l.next()
	}

	id := l.input[start:l.pos]
	return l.tok(TOK_NUMBER, id), nil
}

func (l *Lexer) Next() (Token, error) {
	l.skipWhitespace()
	ch := l.next()

	if ch == EOF {
		return l.tok(EOF, ""), nil
	}
	nextCh := l.peek()

	switch ch {
	case "!":
		if nextCh == "=" {
			l.next()
			return l.tok(TOK_NOT_EQUAL, ch+nextCh), nil
		}
		return l.tok(TOK_BANG, ch), nil
	case "~":
		return l.tok(TOK_TILDE, ch), nil
	case "+":
		return l.tok(TOK_PLUS, ch), nil
	case "-":
		return l.tok(TOK_MINUS, ch), nil
	case "*":
		return l.tok(TOK_ASTERISK, ch), nil
	case "/":
		return l.tok(TOK_SLASH, ch), nil
	case "%":
		return l.tok(TOK_PERCENT, ch), nil
	case "^":
		return l.tok(TOK_CARET, ch), nil
	case "=":
		if nextCh == "=" {
			l.next()
			return l.tok(TOK_DOUBLE_EQUAL, ch+nextCh), nil
		}
		return l.tok(TOK_EQUAL, ch), nil
	case ">":
		if nextCh == "=" {
			l.next()
			return l.tok(TOK_GREATER_THAN_EQUAL, ch+nextCh), nil
		}
		return l.tok(TOK_GREATER_THAN, ch), nil
	case "<":
		if nextCh == "=" {
			l.next()
			return l.tok(TOK_LESS_THAN_EQUAL, ch+nextCh), nil
		}
		return l.tok(TOK_LESS_THAN, ch), nil
	case "&":
		if nextCh == "&" {
			l.next()
			return l.tok(TOK_AND, ch+nextCh), nil
		}
		return Token{}, l.error("Unexpected character: `&` did you mean `&&` ?")
	case "|":
		if nextCh == "|" {
			l.next()
			return l.tok(TOK_OR, ch+nextCh), nil
		}
		return Token{}, l.error("Unexpected character: `|` did you mean `||` ?")
	case "(":
		return l.tok(TOK_LPAREN, ch), nil
	case ")":
		return l.tok(TOK_RPAREN, ch), nil
	case ",":
		return l.tok(TOK_COMMA, ch), nil
	case ".":
		if isValidIdentifierStart(nextCh) {
			return l.readIdentifier()
		}
		if isValidNumber(nextCh) {
			return l.readNumber(nextCh)
		}
		return l.tok(TOK_DOT, ch), nil
	case "#":
		if isValidHex(nextCh) {
			hex, err := l.readHex()
			if err == nil {
				return hex, nil
			}
		}
		if isValidIdentifierStart(nextCh) {
			return l.readIdentifier()
		}
		return l.tok(TOK_HASH, ch), nil
	case "@":
		return l.tok(TOK_AT, ch), nil
	case "{":
		return l.tok(TOK_LSQUIRLY, ch), nil
	case "}":
		return l.tok(TOK_RSQUIRLY, ch), nil
	case ":":
		return l.tok(TOK_COLON, ch), nil
	case ";":
		return l.tok(TOK_SEMICOLON, ch), nil
	case "\"":
		return l.readString(ch)
	case "'":
		return l.readString(ch)
	}

	if isValidIdentifierStart(ch) {
		return l.readIdentifier()
	}

	if isValidNumber(ch) {
		return l.readNumber(ch)
	}

	panic("Not implemented")
}

func main() {
	test := `
"hello" 'world'
		! ~ + - * / % ^ = == != > >= < <= && ||
		( ) , . #
		@ { } : ;
		_main__123
		#id .class
		#id.class
		true false
		!true ~false
	 123 37.5 .22 #22
`
	lexer := NewLexer(test)
	for {
		tok, err := lexer.Next()
		if err != nil {
			panic(err)
		}
		fmt.Println(tok)
		if tok.typ == EOF {
			break
		}
	}
}

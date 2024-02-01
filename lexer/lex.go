package lexer

import "fmt"

type Lexer struct {
	input string
	pos   int
	line  int
	col   int
	done  bool
}

func New(input string) *Lexer {
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
		Typ:   typ,
		Value: value,
		Pos:   l.pos,
		Line:  l.line,
		Col:   l.col - len(value) - 1,
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

	if deci {
		return l.tok(TOK_FLOAT, id), nil
	}
	return l.tok(TOK_INT, id), nil
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
		if nextCh == "-" {
			return l.readIdentifier()
		}
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
		if isValidIdentifierStart(nextCh) {
			return l.readIdentifier()
		}
		return l.tok(TOK_HASH, ch), nil
	case "[":
		return l.tok(TOK_LBRACKET, ch), nil
	case "]":
		return l.tok(TOK_RBRACKET, ch), nil
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
	case "\"", "'":
		return l.readString(ch)
	}

	if isValidIdentifierStart(ch) {
		return l.readIdentifier()
	}

	if isValidNumber(ch) {
		return l.readNumber(ch)
	}

	return Token{}, l.error("Unexpected character: `%s`", ch)
}

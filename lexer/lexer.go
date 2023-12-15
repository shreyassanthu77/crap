package lexer

import (
	"unicode"

	"github.con/shreyascodes-tech/sss-lang/ast"
)

type Lexer struct {
	input string
	pos   int
	line  int
	col   int
}

func New(input string) *Lexer {
	return &Lexer{
		input: input,
		pos:   0,
		line:  1,
		col:   1,
	}
}

func (l *Lexer) Next() Token {
	l.skipWhitespace()
	c := l.next()

	switch c {
	case 0:
		return l.tok(EOF, l.pos-1)
	case '+':
		return l.tok(PLUS, l.pos-1)
	case '-':
		return l.tok(MINUS, l.pos-1)
	case '*':
		return l.tok(ASTERISK, l.pos-1)
	case '/':
		if l.peek() == '*' {
			return l.scanComment()
		}
		return l.tok(SLASH, l.pos-1)
	case '(':
		return l.tok(LPAREN, l.pos-1)
	case ')':
		return l.tok(RPAREN, l.pos-1)
	case '{':
		return l.tok(LSQUIRLY, l.pos-1)
	case '}':
		return l.tok(LSQUIRLY, l.pos-1)
	case '[':
		return l.tok(LBRACKET, l.pos-1)
	case ']':
		return l.tok(RBRACKET, l.pos-1)
	case ':':
		return l.tok(COLON, l.pos-1)
	case ';':
		return l.tok(SEMICOLON, l.pos-1)
	case ',':
		return l.tok(COMMA, l.pos-1)
	case '.':
		if unicode.IsDigit(rune(l.peek())) {
			return l.scanNumber()
		}
		return l.tok(DOT, l.pos-1)
	case '#':
		return l.tok(OCTOTHORPE, l.pos-1)
	case '%':
		return l.tok(PERCENT, l.pos-1)
	case '!':
		return l.tok(EXCLAMATION, l.pos-1)
	case '>':
		return l.tok(GREATER, l.pos-1)
	default:
		if unicode.IsLetter(rune(c)) {
			return l.scanIdentifier()
		} else if unicode.IsDigit(rune(c)) {
			return l.scanNumber()
		} else if c == '"' || c == '\'' {
			return l.scanString(c)
		} else {
			return l.tok(ILLEGAL, l.pos-1)
		}
	}
}

func (l *Lexer) scanComment() Token {
	l.next()
	start := l.pos
	for {
		c := l.next()
		if c == 0 {
			break
		}
		if c == '*' && l.peek() == '/' {
			end := l.pos - 1
			l.next()
			return l.tok(COMMENT, start, end)
		}
	}
	return l.tok(ILLEGAL, start)
}

func (l *Lexer) scanIdentifier() Token {
	start := l.pos - 1
	for {
		c := l.peek()
		if unicode.IsLetter(rune(c)) || unicode.IsDigit(rune(c)) || c == '-' || c == '_' {
			l.next()
		} else {
			break
		}
	}
	return l.tok(IDENT, start)
}

func (l *Lexer) scanNumber() Token {
	start := l.pos - 1
	for {
		c := l.peek()
		if unicode.IsDigit(rune(c)) || c == '.' {
			l.next()
		} else {
			break
		}
	}
	return l.tok(NUMBER, start)
}

func (l *Lexer) scanString(quote byte) Token {
	start := l.pos
	for {
		c := l.next()
		if c == 0 {
			break
		}
		if c == '\\' {
			l.next()
		}
		if c == quote {
			end := l.pos - 1
			return l.tok(STRING, start, end)
		}
	}
	return l.tok(ILLEGAL, start)
}

func (l *Lexer) tok(t TokenType, start int, end ...int) Token {
	if t == EOF {
		return Token{
			Type:  t,
			Value: "",
			Loc:   ast.Location{Line: l.line, Column: l.col, Pos: start, Len: 0},
		}
	}
	length := l.pos - start
	if len(end) > 0 {
		length = end[0] - start
	}
	v := l.input[start : start+length]

	return Token{
		Type:  t,
		Value: v,
		Loc:   ast.Location{Line: l.line, Column: l.col, Pos: start, Len: length},
	}
}

func (l *Lexer) peek() byte {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos]
}

func (l *Lexer) next() byte {
	c := l.peek()
	l.pos++
	if c == '\n' {
		l.line++
		l.col = 1
	} else {
		l.col++
	}
	return c
}

func (l *Lexer) skipWhitespace() {
	for {
		c := l.peek()
		if unicode.IsSpace(rune(c)) {
			l.next()
		} else {
			break
		}
	}
}

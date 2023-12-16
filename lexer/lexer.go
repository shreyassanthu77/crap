package lexer

import (
	"unicode"
)

type Lexer struct {
	input     string
	pos       int
	line      int
	col       int
	startSpan Span
	done      bool
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

func (l *Lexer) Next() Token {
	if l.done {
		return l.token(EOF)
	}

	l.skipWhitespace()

	l.startSpan = l.span()
	c := l.next()

	switch c {
	case 0:
		l.done = true
		return l.token(EOF)
	case '+':
		if n := l.peek(); unicode.IsDigit(rune(n)) {
			l.next() // a number is positive by default
			return l.scanNumber()
		}
		return l.token(PLUS)
	case '-':
		if n := l.peek(); unicode.IsDigit(rune(n)) {
			return l.scanNumber()
		}
		if l.peek() == '-' {
			return l.scanIdentifier()
		}
		return l.token(MINUS)
	case '*':
		return l.token(ASTERISK)
	case '/':
		if l.peek() == '*' {
			return l.scanComment()
		}
		return l.token(SLASH)
	case '=':
		if l.peek() == '=' {
			l.next()
			return l.token(EQUALS)
		}
		return l.token(ASSIGN)
	case '<':
		if l.peek() == '=' {
			l.next()
			return l.token(LESSEQ)
		}
		return l.token(LESS)
	case '>':
		if l.peek() == '=' {
			l.next()
			return l.token(GREATEREQ)
		}
		return l.token(GREATER)
	case '(':
		return l.token(LPAREN)
	case ')':
		return l.token(RPAREN)
	case '{':
		return l.token(LSQUIRLY)
	case '}':
		return l.token(LSQUIRLY)
	case '[':
		return l.token(LBRACKET)
	case ']':
		return l.token(RBRACKET)
	case ':':
		return l.token(COLON)
	case ';':
		return l.token(SEMICOLON)
	case ',':
		return l.token(COMMA)
	case '.':
		if unicode.IsDigit(rune(l.peek())) {
			return l.scanNumber()
		}
		return l.token(DOT)
	case '#':
		if isHexDigit(l.peek()) {
			return l.scanHex()
		}
		return l.token(OCTOTHORPE)
	case '@':
		return l.token(AT)
	case '%':
		return l.token(PERCENT)
	case '!':
		if l.peek() == '=' {
			l.next()
			return l.token(NOTEQUALS)
		}
		return l.token(EXCLAMATION)
	default:
		if unicode.IsLetter(rune(c)) {
			return l.scanIdentifier()
		} else if unicode.IsDigit(rune(c)) {
			return l.scanNumber()
		} else if c == '"' || c == '\'' {
			return l.scanString(c)
		} else {
			return l.token(ILLEGAL)
		}
	}
}

func (l *Lexer) scanComment() Token {
	l.next() // skip the '*'
	start := l.pos
	for {
		c := l.next()
		if c == 0 {
			break
		}
		if c == '*' && l.peek() == '/' {
			end := l.pos - 1
			l.next() // skip the '/'
			return l.token_v(COMMENT, l.input[start:end], l.span())
		}
	}
	return l.token(ILLEGAL)
}

func (l *Lexer) scanIdentifier() Token {
	for {
		c := l.peek()
		if unicode.IsLetter(rune(c)) || unicode.IsDigit(rune(c)) || c == '-' || c == '_' {
			l.next()
		} else {
			break
		}
	}
	return l.token(IDENT)
}

func (l *Lexer) scanNumber() Token {
	for {
		c := l.peek()
		if unicode.IsDigit(rune(c)) || c == '.' {
			l.next()
		} else {
			break
		}
	}
	return l.token(NUMBER)
}

func (l *Lexer) scanHex() Token {
	for {
		c := l.peek()
		if isHexDigit(c) {
			l.next()
		} else {
			break
		}
	}
	return l.token(HEX)
}

func (l *Lexer) scanString(quote byte) Token {
	v := ""
	for {
		c := l.next()
		if c == 0 {
			break
		}
		if c == quote {
			return l.token_v(STRING, v, l.span())
		}
		if c == '\\' {
			switch l.peek() {
			case 'n':
				c = '\n'
			case 'r':
				c = '\r'
			case 't':
				c = '\t'
			case 'f':
				c = '\f'
			case '\\':
				c = '\\'
			case quote:
				c = quote
			default:
				return l.token(ILLEGAL)
			}
			l.next()
		}
		v += string(c)
	}
	return l.token(ILLEGAL)
}

func (l *Lexer) span(offset ...int) Span {
	pos := l.pos
	if len(offset) > 1 {
		panic("too many arguments")
	}
	if len(offset) == 1 {
		pos = offset[0]
	}
	return Span{
		Pos:  pos,
		Line: l.line,
		Col:  l.col,
	}
}

func (l *Lexer) token(t TokenType) Token {
	end := l.span()
	var value string
	if l.pos <= len(l.input) {
		value = l.input[l.startSpan.Pos:l.pos]
	}
	return Token{
		Type:  t,
		Value: value,
		Loc: Location{
			Start: l.startSpan,
			End:   end,
		},
	}
}

func (l *Lexer) token_v(t TokenType, value string, end Span) Token {
	return Token{
		Type:  t,
		Value: value,
		Loc: Location{
			Start: l.startSpan,
			End:   end,
		},
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
	if c == '\n' || c == '\r' || c == '\f' {
		l.line++
		if c == '\r' && l.peek() == '\n' {
			l.pos++
		}
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

package lexer

import (
	"io"
)

const eof = byte(0)

type Lexer struct {
	r    io.Reader
	buf  []byte
	pos  int
	line int
	col  int
	done bool
}

func New(r io.Reader) *Lexer {
	buf := make([]byte, 1024)
	return &Lexer{
		r:   r,
		buf: buf,
	}
}

func (l *Lexer) NextToken() Token {
	panic("TODO")
}

func (l *Lexer) peek(offset int) byte {
	if l.done {
		return eof
	}

	if l.pos+offset >= len(l.buf) {
		n, err := l.r.Read(l.buf[:])
		if err != nil {
			if err == io.EOF {
				l.done = true
				return eof
			}
			panic(err) // TODO: handle error
		}
		l.pos = 0
		l.buf = l.buf[:n]
	}

	if l.pos+offset >= len(l.buf) {
		return eof
	}

	return l.buf[l.pos+offset]
}

func (l *Lexer) next() byte {
	c := l.peek(0)
	if c == eof {
		return eof
	}
	l.pos++
	if c == '\r' && l.peek(1) == '\n' {
		l.pos++ // skip \r if followed by \n
	}
	if c == '\n' {
		l.line++
		l.col = 0
	} else {
		l.col++
	}
	return c
}

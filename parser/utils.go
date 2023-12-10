package parser

import "fmt"

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isAlpha(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')
}

func isAlphaNumeric(ch byte) bool {
	return isAlpha(ch) || isDigit(ch)
}

func (p *Parser) logState(msg string, args ...interface{}) {
	fmt.Printf(
		`
===============
%s
---------------
pos: %d, line: %d, col: %d
%s_%c_%s
===============
`,
		fmt.Sprintf(msg, args...),
		p.pos, p.line, p.col,
		p.input[max(0, p.pos-10):p.pos],
		p.input[p.pos],
		p.input[p.pos+1:min(p.pos+10, p.length)],
	)
}

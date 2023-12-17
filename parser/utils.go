package parser

import "github.con/shreyascodes-tech/sss-lang/lexer"

func isCombinator(t lexer.Token) bool {
	return t.Type == lexer.PLUS ||
		t.Type == lexer.GREATER
}

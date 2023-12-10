package lexer

import "github.con/shreyascodes-tech/sss-lang/ast"

type Token struct {
	Type  TokenType
	Value string
	Loc   ast.Location
}

type TokenType int

// CSS Tokens
const (
	ILLEGAL TokenType = iota
	EOF
	WS
	COMMENT

	IDENT
	NUMBER
	STRING

	COLON
	SEMICOLON
	COMMA
	LSQUIRLY
	RSQUIRLY
	LPAREN
	RPAREN
	LBRACKET
	RBRACKET
	DOT
	OCTOTHORPE
	PLUS
	MINUS
	ASTERISK
	SLASH
	PERCENT
	EXCLAMATION
	GREATER
)

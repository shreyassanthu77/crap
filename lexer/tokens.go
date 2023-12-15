package lexer

import (
	"fmt"

	"github.con/shreyascodes-tech/sss-lang/ast"
)

type Token struct {
	Type  TokenType
	Value string
	Loc   ast.Location
}

type TokenType string

// CSS Tokens
const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF               = "EOF"
	COMMENT           = "COMMENT"

	IDENT  = "IDENT"
	NUMBER = "NUMBER"
	STRING = "STRING"

	COLON       = "COLON"
	SEMICOLON   = "SEMICOLON"
	COMMA       = "COMMA"
	LSQUIRLY    = "LSQUIRLY"
	RSQUIRLY    = "RSQUIRLY"
	LPAREN      = "LPAREN"
	RPAREN      = "RPAREN"
	LBRACKET    = "LBRACKET"
	RBRACKET    = "RBRACKET"
	DOT         = "DOT"
	OCTOTHORPE  = "OCTOTHORPE"
	PLUS        = "PLUS"
	MINUS       = "MINUS"
	ASTERISK    = "ASTERISK"
	SLASH       = "SLASH"
	PERCENT     = "PERCENT"
	EXCLAMATION = "EXCLAMATION"
	GREATER     = "GREATER"
)

func (t Token) String() string {
	if t.Type == EOF {
		return "--EOF--"
	}
	fmtStr := "Token<%s> /%s/ %d-%d %d:%d"
	return fmt.Sprintf(fmtStr, t.Type, t.Value, t.Loc.Pos, t.Loc.Len, t.Loc.Line, t.Loc.Column)
}

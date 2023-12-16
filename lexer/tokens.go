package lexer

import (
	"fmt"
)

type Span struct {
	Pos  int
	Line int
	Col  int
}

type Location struct {
	Start Span
	End   Span
}

func (l Location) String() string {
	return fmt.Sprintf("%d:%d-%d:%d", l.Start.Line, l.Start.Col, l.End.Line, l.End.Col)
}

type Token struct {
	Type  TokenType
	Value string
	Loc   Location
}

type TokenType string

// CSS Tokens
const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF               = "EOF"
	COMMENT           = "COMMENT"

	IDENT  = "IDENT"
	NUMBER = "NUMBER"
	HEX    = "HEX"
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
	AT          = "AT"
	PLUS        = "PLUS"
	MINUS       = "MINUS"
	ASTERISK    = "ASTERISK"
	SLASH       = "SLASH"
	PERCENT     = "PERCENT"
	ASSIGN      = "ASSIGN"
	EQUALS      = "EQUALS"
	NOTEQUALS   = "NOTEQUALS"
	EXCLAMATION = "EXCLAMATION"
	GREATER     = "GREATER"
	GREATEREQ   = "GREATEREQ"
	LESS        = "LESS"
	LESSEQ      = "LESSEQ"
)

func (t Token) String() string {
	if t.Type == EOF {
		return "--EOF--"
	}
	fmtStr := "Token<%s> /%s/ %s"
	return fmt.Sprintf(fmtStr, t.Type, t.Value, t.Loc)
}

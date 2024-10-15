//go:generate go run golang.org/x/tools/cmd/stringer -type=TokenType
package token

import (
	"fmt"
)

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF

	// Identifiers + literals
	IDENT
	NUMBER
	STRING

	// Operators
	ASSIGN
	PLUS
	MINUS
	BANG
	ASTERISK
	SLASH
	EQ
	NOT_EQ

	LT
	GT

	AND
	OR

	// Delimiters
	COMMA
	DOT
	SEMICOLON
	COLON

	LPAREN
	RPAREN
	LBRACE
	RBRACE
	LBRACKET
	RBRACKET

	// Keywords
	FUNCTION
	EACH
	LET
	TRUE
	FALSE
	IF
	ELSE
	RETURN
	PRINT
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"and":    AND,
	"or":     OR,
	"each":   EACH,
	"print":  PRINT,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENT
}

type Pos struct {
	Start, End int
}
type Token struct {
	Pos    Pos       // offset from the start of the source file
	Type   TokenType // the number of characters from the start of the line to the start of the token
	Lexeme string    // the string that was parsed as this token
}

func NewToken(tokenType TokenType, lexeme string, pos Pos) Token {
	return Token{
		Pos:    pos,
		Type:   tokenType,
		Lexeme: lexeme,
	}
}

func (t *Token) String() string {
	return fmt.Sprintf("[offset: %d] %s: %s", t.Pos, t.Type, t.Lexeme)
}

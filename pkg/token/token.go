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
	IMPORT
	EACH
	WHILE
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
	"import": IMPORT,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"and":    AND,
	"or":     OR,
	"each":   EACH,
	"while":  WHILE,
	"print":  PRINT,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENT
}

type Pos struct {
	Src        *string
	Start, End int
}

// returns the line:column pair for the token
func (p *Pos) Position() (int, int) {
	col := 1
	line := 1

	position := 0
	// increment col for each loop.
	// when we see a '\n', reset col and keep moving
	// return col at end
	for position < p.Start {
		if (*p.Src)[position] == '\n' {
			line += 1
			col = 1
		} else {
			col += 1
		}

		position += 1
	}

	return line, col
}

func (p Pos) String() string {
	line, col := p.Position()
	return fmt.Sprintf("[%d:%d]", line, col)
}

type Token struct {
	Pos    Pos       // The postion of the lexeme in the source
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
	return fmt.Sprintf("%s %s: %s", t.Pos, t.Type, t.Lexeme)
}

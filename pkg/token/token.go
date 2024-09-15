//go:generate go run golang.org/x/tools/cmd/stringer -type=TokenType
package token

import "fmt"

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF

	// Identifiers + literals
	IDENT
	INT
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
	LET
	TRUE
	FALSE
	IF
	ELSE
	RETURN
)

type Token struct {
	Line      int       // what line the token i parsed from the input text
	TokenType TokenType // the number of characters from the start of the line to the start of the token
	Lexeme    string    // the string that was parsed as this token
	Literal   any       // the literal value of the token. int/string/bool/nil
}

func NewToken(tokenType TokenType, lexeme string, literal any, line, col int) Token {
	return Token{
		Line:      line,
		TokenType: tokenType,
		Lexeme:    lexeme,
		Literal:   literal,
	}
}

func (t *Token) String() string {
	return fmt.Sprintf("[%d] %s: %s", t.Line, t.TokenType, t.Lexeme)
}

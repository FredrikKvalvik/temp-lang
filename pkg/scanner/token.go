//go:generate go run golang.org/x/tools/cmd/stringer -type=TokenType
package scanner

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
	// what line the token i parsed from the input text
	Line int
	// the number of characters from the start of the line to the start of the token
	Col int
	// the type of token
	TokenType TokenType
	// the string that was parsed as this token
	Lexeme string
	// the literal value of the token. int/string/bool/nil
	Literal any
}

func NewToken(tokenType TokenType, lexeme string, literal any, line, col int) Token {
	return Token{
		Line:      line,
		Col:       col,
		TokenType: tokenType,
		Lexeme:    lexeme,
		Literal:   literal,
	}
}

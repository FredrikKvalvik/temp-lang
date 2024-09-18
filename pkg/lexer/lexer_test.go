package lexer

import (
	"testing"

	"github.com/fredrikkvalvik/temp-lang/pkg/token"
)

func TestNextTokoen(t *testing.T) {
	input := `let five = 5
let ten = 10;
	let add = fn(x, y) {
     x + y;
};

// comment

let result = add(five, ten);
  !-/*5;
5 < 10 > 5;

if (5 < 10) {
   return true;
} else {
   return false;
}

10 == 10 // comment after expression
10 != 9
true and false
true or false
"foobar"
"foo bar"
[1, 2];
{"foo": "bar"}
`

	tests := []struct {
		expectedType    token.TokenType
		exptectedLexeme string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.NUMBER, "5"},
		{token.SEMICOLON, "\n"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.NUMBER, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.NUMBER, "5"},
		{token.SEMICOLON, ";"},
		{token.NUMBER, "5"},
		{token.LT, "<"},
		{token.NUMBER, "10"},
		{token.GT, ">"},
		{token.NUMBER, "5"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.NUMBER, "5"},
		{token.LT, "<"},
		{token.NUMBER, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.NUMBER, "10"},
		{token.EQ, "=="},
		{token.NUMBER, "10"},
		{token.SEMICOLON, "\n"},
		{token.NUMBER, "10"},
		{token.NOT_EQ, "!="},
		{token.NUMBER, "9"},
		{token.SEMICOLON, "\n"},
		{token.TRUE, "true"},
		{token.AND, "and"},
		{token.FALSE, "false"},
		{token.SEMICOLON, "\n"},
		{token.TRUE, "true"},
		{token.OR, "or"},
		{token.FALSE, "false"},
		{token.SEMICOLON, "\n"},
		{token.STRING, "\"foobar\""},
		{token.SEMICOLON, "\n"},
		{token.STRING, "\"foo bar\""},
		{token.SEMICOLON, "\n"},
		{token.LBRACKET, "["},
		{token.NUMBER, "1"},
		{token.COMMA, ","},
		{token.NUMBER, "2"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},
		{token.LBRACE, "{"},
		{token.STRING, "\"foo\""},
		{token.COLON, ":"},
		{token.STRING, "\"bar\""},
		{token.RBRACE, "}"},
		{token.EOF, string(byte(0))},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q. lexeme=%s, prevToken=%s",
				i, tt.expectedType, tok.Type, tok.Lexeme, l.previousToken.Type)
		}

		if tok.Lexeme != tt.exptectedLexeme {
			t.Fatalf("tests[%d] - lexeme wrong. expected=%q, got=%q",
				i, tt.exptectedLexeme, tok.Lexeme)
		}
	}
}

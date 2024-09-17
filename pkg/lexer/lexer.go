package lexer

import (
	"fmt"
	"strconv"

	"github.com/fredrikkvalvik/temp-lang/pkg/token"
)

type Lexer struct {
	source string
	tokens []token.Token

	position     int  // position of the current lexeme
	readPosition int  // nexh char to be read
	ch           byte // current ch
	line         int  // current line in source

	errors []error
}

func New(source string) *Lexer {
	l := &Lexer{

		source: source,
		tokens: []token.Token{},

		line: 1,

		errors: make([]error, 0),
	}
	l.advance()
	return l
}

func (l *Lexer) Errors() []error {
	return l.errors
}

func (l *Lexer) DidError() bool {
	return len(l.errors) > 0
}

// pull tokens when needed
func (l *Lexer) NextToken() token.Token {
	return l.scanToken()
}

func (l *Lexer) scanToken() token.Token {
	var tok token.Token

	l.whitespace()

	switch l.ch {
	case '(':
		tok = l.getToken(token.LPAREN, string(l.ch), nil)
	case ')':
		tok = l.getToken(token.RPAREN, string(l.ch), nil)
	case '{':
		tok = l.getToken(token.LBRACE, string(l.ch), nil)
	case '}':
		tok = l.getToken(token.RBRACE, string(l.ch), nil)
	case '[':
		tok = l.getToken(token.LBRACKET, string(l.ch), nil)
	case ']':
		tok = l.getToken(token.RBRACKET, string(l.ch), nil)
	case '+':
		tok = l.getToken(token.PLUS, string(l.ch), nil)
	case '-':
		tok = l.getToken(token.MINUS, string(l.ch), nil)
	case '*':
		tok = l.getToken(token.ASTERISK, string(l.ch), nil)
	case ',':
		tok = l.getToken(token.COMMA, string(l.ch), nil)
	case '.':
		tok = l.getToken(token.DOT, string(l.ch), nil)
	case ':':
		tok = l.getToken(token.COLON, string(l.ch), nil)
	case ';':
		tok = l.getToken(token.SEMICOLON, string(l.ch), nil)
	case '<':
		tok = l.getToken(token.LT, string(l.ch), nil)
	case '>':
		tok = l.getToken(token.GT, string(l.ch), nil)

	case '/':

		tok = l.getToken(token.SLASH, string(l.ch), nil)

	case '=':
		if l.peek() == '=' {
			l.advance()
			tok = l.getToken(token.EQ, "==", nil)
		} else {
			tok = l.getToken(token.ASSIGN, string(l.ch), nil)
		}

	case '!':
		if l.peek() == '=' {
			ch := l.ch
			l.advance()
			tok = l.getToken(token.NOT_EQ, string(ch)+string(l.ch), nil)
		} else {
			tok = l.getToken(token.BANG, string(l.ch), nil)
		}

	case '"':
		tokType := token.STRING
		lexeme, literal := l.readString()
		tok = l.getToken(tokType, lexeme, literal)

	case 0:
		tok = l.getToken(token.EOF, string(l.ch), nil)

	default:
		if isLetter(l.ch) {
			// can be identier or reserved word
			lexeme := l.readIdentifier()
			tokType := token.LookupIdent(lexeme)
			tok = l.getToken(tokType, lexeme, nil)

		} else if isDigit(l.ch) {
			lexeme, literal, err := l.readNumber()
			if err != nil {
				tok = l.getToken(token.ILLEGAL, "", nil)
				l.error(err)
				break
			}
			tok = l.getToken(token.NUMBER, lexeme, literal)

		} else {
			tok = l.getToken(token.ILLEGAL, string(l.ch), nil)
			l.error(fmt.Errorf("Unexpected character"))
		}
	}

	l.advance()
	return tok
}

func (l *Lexer) getToken(t token.TokenType, lexeme string, literal any) token.Token {
	tok := token.NewToken(t, lexeme, literal, l.line, l.readPosition)

	return tok
}

func (l *Lexer) advance() {
	if l.atEnd() {
		l.ch = 0
	} else {
		l.ch = l.source[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peek() byte {
	if l.atEnd() {
		return 0
	}

	return l.source[l.readPosition]
}

func (l *Lexer) atEnd() bool {
	return l.readPosition >= len(l.source)
}

func (l *Lexer) whitespace() {
	for {
		if l.atEnd() {
			break
		}

		if l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
			l.advance()
			continue
		}
		if l.ch == '\n' {
			l.advance()
			l.line += 1
			continue
		}

		break
	}
}

// returns the lexeme and the literal value of the string
func (l *Lexer) readString() (string, string) {
	for {
		if l.peek() != '"' {
			l.readPosition += 1
		} else {
			// consume the ending quote
			l.readPosition += 1
			break
		}
	}
	lexeme := l.source[l.position:l.readPosition]
	literal := lexeme[1 : len(lexeme)-1]

	return lexeme, literal
}

func (l *Lexer) readNumber() (string, float64, error) {
	for {
		if !isDigit(l.peek()) {
			break
		}

		l.readPosition += 1
		continue
	}
	if l.peek() == '.' {
		// consume .
		l.readPosition += 1

		// if the next char is not a number, then the token is invalid
		if isDigit(l.peek()) {
			return "", 0, nil
		}
		// parse decimal digits
		for {
			if !isDigit(l.peek()) {
				break
			}

			l.readPosition += 1
			continue
		}
	}

	lexeme := l.source[l.position:l.readPosition]
	literal, err := strconv.ParseFloat(lexeme, 64)
	if err != nil {
		return "", 0, err
	}

	return lexeme, literal, nil
}

func (l *Lexer) readIdentifier() string {
	for !l.atEnd() && (isLetter(l.peek()) || isDigit(l.peek())) {
		l.readPosition += 1
	}

	return l.source[l.position:l.readPosition]
}

func (l *Lexer) error(err error) {
	l.errors = append(l.errors, err)
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

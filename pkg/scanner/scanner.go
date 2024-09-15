package scanner

import (
	"fmt"

	"github.com/fredrikkvalvik/temp-lang/pkg/token"
)

type Scanner struct {
	source string
	tokens []token.Token

	position     int  // position of the current lexeme
	readPosition int  // nexh char to be read
	ch           byte // current ch
	line         int  // current line in source

	errors []error
}

func New(source string) *Scanner {
	s := &Scanner{

		source: source,
		tokens: []token.Token{},

		line: 1,

		errors: make([]error, 0),
	}
	s.advance()
	return s
}

func (s *Scanner) Errors() []error {
	return s.errors
}

func (s *Scanner) DidError() bool {
	return len(s.errors) > 0
}

// pull tokens when needed
func (s *Scanner) NextToken() token.Token {
	return s.scanToken()
}

func (s *Scanner) scanToken() token.Token {
	var tok token.Token

	s.whitespace()

	switch s.ch {
	case '(':
		tok = s.getToken(token.LPAREN, string(s.ch), nil)
	case ')':
		tok = s.getToken(token.RPAREN, string(s.ch), nil)
	case '{':
		tok = s.getToken(token.LBRACE, string(s.ch), nil)
	case '}':
		tok = s.getToken(token.RBRACE, string(s.ch), nil)
	case '!':
		tok = s.getToken(token.BANG, string(s.ch), nil)
	case '+':
		tok = s.getToken(token.PLUS, string(s.ch), nil)
	case '-':
		tok = s.getToken(token.MINUS, string(s.ch), nil)
	case '*':
		tok = s.getToken(token.ASTERISK, string(s.ch), nil)
	case ',':
		tok = s.getToken(token.COMMA, string(s.ch), nil)
	case '.':
		tok = s.getToken(token.DOT, string(s.ch), nil)
	case ':':
		tok = s.getToken(token.COLON, string(s.ch), nil)
	case ';':
		tok = s.getToken(token.SEMICOLON, string(s.ch), nil)

	case '/':
		tok = s.getToken(token.SLASH, string(s.ch), nil)

	case '<':
		tok = s.getToken(token.LT, string(s.ch), nil)
	case '>':
		tok = s.getToken(token.GT, string(s.ch), nil)

	case 0:
		tok = s.getToken(token.EOF, string(s.ch), nil)

	default:
		tok = s.getToken(token.ILLEGAL, string(s.ch), nil)
		s.error(fmt.Errorf("Unexpected character"))
	}

	s.advance()
	return tok
}

func (s *Scanner) getToken(t token.TokenType, lexeme string, literal any) token.Token {
	tok := token.NewToken(t, string(s.ch), literal, s.line, s.readPosition)

	return tok
}

func (s *Scanner) advance() {
	if s.atEnd() {
		s.ch = 0
	} else {
		s.ch = s.source[s.readPosition]
	}

	s.position = s.readPosition
	s.readPosition += 1
}

func (s *Scanner) whitespace() {
	for {
		if s.atEnd() {
			break
		}

		if s.ch == ' ' || s.ch == '\t' || s.ch == '\r' {
			s.advance()
			continue
		}
		if s.ch == '\n' {
			s.advance()
			s.line += 1
			continue
		}

		break
	}
}

func (s *Scanner) atEnd() bool {
	return s.readPosition >= len(s.source)
}

func (s *Scanner) error(err error) {
	s.errors = append(s.errors, err)
}

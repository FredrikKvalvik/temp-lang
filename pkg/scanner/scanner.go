package scanner

import (
	"fmt"
	"strconv"

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
	case '[':
		tok = s.getToken(token.LBRACKET, string(s.ch), nil)
	case ']':
		tok = s.getToken(token.RBRACKET, string(s.ch), nil)
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
	case '<':
		tok = s.getToken(token.LT, string(s.ch), nil)
	case '>':
		tok = s.getToken(token.GT, string(s.ch), nil)

	case '/':

		tok = s.getToken(token.SLASH, string(s.ch), nil)

	case '=':
		if s.peek() == '=' {
			s.advance()
			tok = s.getToken(token.EQ, "==", nil)
		} else {
			tok = s.getToken(token.ASSIGN, string(s.ch), nil)
		}

	case '!':
		if s.peek() == '=' {
			ch := s.ch
			s.advance()
			tok = s.getToken(token.NOT_EQ, string(ch)+string(s.ch), nil)
		} else {
			tok = s.getToken(token.BANG, string(s.ch), nil)
		}

	case '"':
		tokType := token.STRING
		lexeme, literal := s.readString()
		tok = s.getToken(tokType, lexeme, literal)

	case 0:
		tok = s.getToken(token.EOF, string(s.ch), nil)

	default:
		if isLetter(s.ch) {
			// can be identier or reserved word
			lexeme := s.readIdentifier()
			tokType := token.LookupIdent(lexeme)
			tok = s.getToken(tokType, lexeme, nil)

		} else if isDigit(s.ch) {
			lexeme, literal, err := s.readNumber()
			if err != nil {
				tok = s.getToken(token.ILLEGAL, "", nil)
				s.error(err)
				break
			}
			tok = s.getToken(token.NUMBER, lexeme, literal)

		} else {
			tok = s.getToken(token.ILLEGAL, string(s.ch), nil)
			s.error(fmt.Errorf("Unexpected character"))
		}
	}

	s.advance()
	return tok
}

func (s *Scanner) getToken(t token.TokenType, lexeme string, literal any) token.Token {
	tok := token.NewToken(t, lexeme, literal, s.line, s.readPosition)

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

func (s *Scanner) peek() byte {
	if s.atEnd() {
		return 0
	}

	return s.source[s.readPosition]
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

// returns the lexeme and the literal value of the string
func (s *Scanner) readString() (string, string) {
	for {
		if s.peek() != '"' {
			s.readPosition += 1
		} else {
			// consume the ending quote
			s.readPosition += 1
			break
		}
	}
	lexeme := s.source[s.position:s.readPosition]
	literal := lexeme[1 : len(lexeme)-1]

	return lexeme, literal
}

func (s *Scanner) readNumber() (string, float64, error) {
	for {
		if !isDigit(s.peek()) {
			break
		}

		s.readPosition += 1
		continue
	}
	if s.peek() == '.' {
		// consume .
		s.readPosition += 1

		// if the next char is not a number, then the token is invalid
		if isDigit(s.peek()) {
			return "", 0, nil
		}
		// parse decimal digits
		for {
			if !isDigit(s.peek()) {
				break
			}

			s.readPosition += 1
			continue
		}
	}

	lexeme := s.source[s.position:s.readPosition]
	literal, err := strconv.ParseFloat(lexeme, 64)
	if err != nil {
		return "", 0, err
	}

	return lexeme, literal, nil
}

func (s *Scanner) readIdentifier() string {
	for !s.atEnd() && (isLetter(s.peek()) || isDigit(s.peek())) {
		s.readPosition += 1
	}

	return s.source[s.position:s.readPosition]
}

func (s *Scanner) atEnd() bool {
	return s.readPosition >= len(s.source)
}

func (s *Scanner) error(err error) {
	s.errors = append(s.errors, err)
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

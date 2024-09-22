package lexer

import (
	"fmt"
	"strconv"

	"github.com/fredrikkvalvik/temp-lang/pkg/token"
)

// TODO: add implicit semi-colon with automatic insertion
// ref: https://go.dev/doc/effective_go#semicolons

type Lexer struct {
	source string
	tokens []token.Token

	position     int  // position of the current lexeme
	readPosition int  // nexh char to be read
	ch           byte // current ch
	line         int  // current line in source

	previousToken *token.Token

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

// used for error messages
func (l *Lexer) GetTokenColumn(tok *token.Token) int {
	offset := tok.Offset

	return l.getColFromOffset(offset)
}

func (l *Lexer) getColFromOffset(offset int) int {
	col := 1

	position := 0
	// increment col for each loop.
	// when we see a '\n', reset col and keep moving
	// return col at end
	for position < offset {
		if l.source[position] == '\n' {
			col = 1
		} else {
			col += 1
		}
		position += 1
	}

	return col
}

// pull tokens when needed
func (l *Lexer) NextToken() token.Token {
	tok := l.scanToken()
	l.previousToken = &tok
	return tok
}

func (l *Lexer) scanToken() token.Token {
	var tok token.Token

	// redo label for comments. instead of doing a new loop, we go back to the top and start from the top again
REDO:
	// might return semicolon token
	if terminal := l.whitespace(); terminal != nil {
		return *terminal
	}

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
		if l.peek() == '/' {
			for l.ch != '\n' {
				l.advance()
			}
			// use goto to jump back to the top to parse next token.
			// this also works with the whitespace call at the start of the function. It will check the newline and see if
			// it satisfies automatic ';' insertion
			goto REDO

		} else {
			tok = l.getToken(token.SLASH, string(l.ch), nil)
		}

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
		if l.newlineIsTerminal() {
			tok = l.getToken(token.SEMICOLON, string(l.ch), nil)
		} else {
			tok = l.getToken(token.EOF, string(l.ch), nil)
		}

	default:
		if isLetter(l.ch) {
			// can be identier or reserved word
			lexeme := l.readIdentifier()
			tokType := token.LookupIdent(lexeme)
			tok = l.getToken(tokType, lexeme, nil)

		} else if isDigit(l.ch) {
			lexeme, literal, err := l.readNumber()
			if err != nil {
				tok = l.getToken(token.ILLEGAL, lexeme, nil)
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
	tok := token.Token{
		Type:    t,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    l.line,
		Offset:  l.position,
	}

	return tok
}

func (l *Lexer) advance() {
	if l.atEnd() {
		l.ch = 0
	} else {
		l.ch = l.source[l.readPosition]
		l.position = l.readPosition
		l.readPosition += 1
	}
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

func (l *Lexer) whitespace() *token.Token {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == '\n' {
		if l.ch == '\n' {
			if t := l.newline(); t != nil {
				return t
			}
		} else {
			l.advance()
		}
	}

	return nil
}

// handles newline for linenumber and returns a semicolon if the conditions are correct
func (l *Lexer) newline() *token.Token {
	if l.newlineIsTerminal() {
		tok := l.getToken(token.SEMICOLON, string(l.ch), nil)
		l.line += 1
		l.advance()
		return &tok

	} else {
		l.advance()
		l.line += 1
		return nil
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

// returns the lexeme and literal of string.
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
		if !isDigit(l.peek()) {
			err := fmt.Errorf("[%d:%d]: expected digit, got=%s\n",
				l.line,
				l.getColFromOffset(l.readPosition-1),
				string(l.peek()))
			return "", 0, err

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

	literal, _ := strconv.ParseFloat(lexeme, 64)

	return lexeme, literal, nil
}

func (l *Lexer) readIdentifier() string {
	for !l.atEnd() && (isLetter(l.peek()) || isDigit(l.peek())) {
		l.readPosition += 1
	}

	return l.source[l.position:l.readPosition]
}

// returns true if the previous token followed by a new line satisifies automatic insertion of semicolon
func (l *Lexer) newlineIsTerminal() bool {
	if l.previousToken == nil {
		return false
	}

	switch l.previousToken.Type {
	case token.IDENT:
		return true
	case token.STRING:
		return true
	case token.NUMBER:
		return true
	case token.RPAREN:
		return true
	case token.RBRACKET:
		return true
	case token.FALSE:
		return true
	case token.TRUE:
		return true
	case token.RETURN:
		return true
	}

	return false
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

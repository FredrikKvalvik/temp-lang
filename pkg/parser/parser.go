package parser

import (
	"github.com/fredrikkvalvik/temp-lang/pkg/ast"
	"github.com/fredrikkvalvik/temp-lang/pkg/lexer"
	"github.com/fredrikkvalvik/temp-lang/pkg/token"
)

const (
	_ int = iota
	LOWEST
	ASSIGN      // =
	OR          // or
	AND         // and
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         //+ -
	PRODUCT     //* /
	PREFIX      //-X or !X
	CALL        // myFunction(X)
	INDEX       // array[index]
)

var stickinessMap = map[token.TokenType]int{
	token.OR:       OR,
	token.AND:      AND,
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	infixParselets  map[token.TokenType]infixFn
	prefixParselets map[token.TokenType]prefixFn

	errors []error
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,

		errors: make([]error, 0),

		infixParselets:  map[token.TokenType]infixFn{},
		prefixParselets: map[token.TokenType]prefixFn{},
	}

	// prepare curToken and peekToken
	p.advance()
	p.advance()

	p.registerPrefix(token.NUMBER, p.parseNumberLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(token.FALSE, p.parseBooleanLiteral)

	p.registerPrefix(token.BANG, p.parsePrefix)
	p.registerPrefix(token.MINUS, p.parsePrefix)

	return p
}

func (p *Parser) DidError() bool {
	return len(p.errors) > 0
}

func (p *Parser) Errors() []error {
	return p.errors
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = make([]ast.Stmt, 0)

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.advance()
	}

	return program
}

func (p *Parser) peekStickiness() int {
	if s, ok := stickinessMap[p.peekToken.Type]; ok {
		return s
	}
	return LOWEST
}

func (p *Parser) curStickiness() int {
	if s, ok := stickinessMap[p.curToken.Type]; ok {
		return s
	}
	return LOWEST
}

// consume current token
func (p *Parser) advance() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) atEnd() bool {
	return p.curToken.Type == token.EOF
}

func (p *Parser) curTokenIs(typ token.TokenType) bool {
	return p.curToken.Type == typ
}

func (p *Parser) peekTokenIs(typ token.TokenType) bool {
	return p.peekToken.Type == typ
}

// consume cur token if cur token == typ
func (p *Parser) expect(typ token.TokenType) bool {
	if p.curToken.Type == typ {
		p.advance()
		return true
	}

	p.expectCurError(typ)
	return false
}

// consume cur token if peek token == typ
func (p *Parser) expectPeek(typ token.TokenType) bool {
	if p.peekToken.Type == typ {
		p.advance()
		return true
	}

	p.expectPeekError(typ)
	return false
}

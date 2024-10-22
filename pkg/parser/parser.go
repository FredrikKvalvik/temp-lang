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
	GET         // obj.property
	GROUPING    // ( expression )
	INDEX       // array[index]
)

var stickinessMap = map[token.TokenType]int{
	token.ASSIGN:   ASSIGN,
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
	token.FUNCTION: CALL,
	token.DOT:      GET,
	token.LPAREN:   GROUPING,
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

type MyString = *string

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

	// literals
	p.registerPrefix(token.NUMBER, p.parseNumberLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(token.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(token.IDENT, p.parseIdent)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.LBRACKET, p.parseListLiteralExpression)
	p.registerPrefix(token.LBRACE, p.parseMapLiteralExpression)

	// complex literals
	p.registerInfix(token.LPAREN, p.parseCall)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)

	// prefix
	p.registerPrefix(token.BANG, p.parsePrefix)
	p.registerPrefix(token.MINUS, p.parsePrefix)
	p.registerPrefix(token.LPAREN, p.parseParenPrefix)

	// Infix
	p.registerInfix(token.PLUS, p.parseBinary)
	p.registerInfix(token.MINUS, p.parseBinary)
	p.registerInfix(token.ASTERISK, p.parseBinary)
	p.registerInfix(token.SLASH, p.parseBinary)
	p.registerInfix(token.LT, p.parseBinary)
	p.registerInfix(token.GT, p.parseBinary)
	p.registerInfix(token.EQ, p.parseBinary)
	p.registerInfix(token.NOT_EQ, p.parseBinary)

	// get
	p.registerInfix(token.DOT, p.parseGetExpression)

	// logical
	p.registerInfix(token.AND, p.parseLogical)
	p.registerInfix(token.OR, p.parseLogical)

	// assign
	p.registerInfix(token.ASSIGN, p.parseAssign)

	return p
}

func (p *Parser) DidError() bool {
	return len(p.errors) > 0
}

func (p *Parser) Errors() []error {
	return p.errors
}

func (p *Parser) GetErrorPosition(tok token.Token) (int, int) {
	line, col := p.l.GetTokenPosition(&tok)

	return line, col
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = make([]ast.Stmt, 0)

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		} else {
			p.recover()
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

// will consume the current token if peekToken == acceptedToken
func (p *Parser) consume(acceptedToken token.TokenType) {
	if p.peekTokenIs(acceptedToken) {
		p.advance()
	}
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

// TODO: add error message to expect if it fails
//
// advances if the typ == peekToken.Type and return true
// if not, return false and stay
func (p *Parser) expectPeek(typ token.TokenType) bool {

	if p.peekTokenIs(typ) {
		p.advance()
		return true
	}

	p.expectPeekError(typ)
	return false
}
func (p *Parser) expectCur(typ token.TokenType) bool {

	if p.curTokenIs(typ) {
		p.advance()
		return true
	}

	p.expectPeekError(typ)
	return false
}

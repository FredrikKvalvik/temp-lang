package parser

import (
	"fmt"

	"github.com/fredrikkvalvik/temp-lang/pkg/ast"
	"github.com/fredrikkvalvik/temp-lang/pkg/token"
)

func (p *Parser) parseStatement() ast.Stmt {

	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatment()

	default:
		return nil
	}

}

func (p *Parser) parseLetStatment() *ast.LetStmt {
	letStmt := &ast.LetStmt{Token: p.curToken}

	// standing on "let"
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	// standing on "ident"
	letStmt.Name = &ast.IdentifierExpr{
		Token: p.curToken,
		Value: p.curToken.Lexeme,
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// standing on assign
	letStmt.Value = p.parseExpression(LOWEST)

	// standing on "expression"
	if !p.expectPeek(token.SEMICOLON) {
		p.errors = append(p.errors, fmt.Errorf("expected ';', got=%s\n", p.peekToken.Type))
	}

	return letStmt
}

// advances one token and tries to parse an expression based on curToken
func (p *Parser) parseExpression(precedence int) ast.Expr {
	p.advance()

	prefix := p.prefixParselets[p.curToken.Type]

	if prefix == nil {
		p.errors = append(p.errors, fmt.Errorf("No prefix parslet for tokenType=%s\n", p.curToken.Type))
		return nil
	}

	left := prefix()

	return left
}

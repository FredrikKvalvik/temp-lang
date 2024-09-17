package parser

import (
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

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	ident := &ast.IdentifierExpr{
		Token: p.curToken,
		Value: p.curToken.Lexeme,
	}
	letStmt.Name = ident

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: skipping expression for now
	for p.curTokenIs(token.SEMICOLON) {
		p.advance()
	}

	return letStmt
}

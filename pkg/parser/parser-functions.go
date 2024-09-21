package parser

import (
	"github.com/fredrikkvalvik/temp-lang/pkg/ast"
	"github.com/fredrikkvalvik/temp-lang/pkg/token"
)

func (p *Parser) parseStatement() ast.Stmt {

	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatment()
	case token.IF:
		return p.parseIfStatement()

	default:
		return p.parseExpressionStatement()
	}

}

func (p *Parser) parseLetStatment() *ast.LetStmt {
	// let   ident    =    "hei"
	// ^
	letStmt := &ast.LetStmt{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	// let   ident    =    "hei"
	//       ^
	letStmt.Name = &ast.IdentifierExpr{
		Token: p.curToken,
		Value: p.curToken.Lexeme,
	}

	if !p.expectPeek(token.ASSIGN) {
		p.recover()
		return nil
	}
	// let   ident    =    "hei"
	//                ^

	// consume '=' to prepare parseExpression
	p.advance()

	// let   ident    =    "hei"
	//                     ^
	letStmt.Value = p.parseExpression(LOWEST)

	// standing at end of "expression"
	if !p.expectPeek(token.SEMICOLON) {
		p.recover()
		return nil
	}

	return letStmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStmt {
	block := &ast.BlockStmt{
		Token:      p.curToken,
		Statements: make([]ast.Stmt, 0),
	}

	// consunme '{'
	p.advance()

	for !p.curTokenIs(token.RBRACE) {
		stmt := p.parseStatement()
		block.Statements = append(block.Statements, stmt)
	}

	return block
}

func (p *Parser) parseIfStatement() *ast.IfStmt {
	ifstmt := &ast.IfStmt{
		Token: p.curToken,
	}

	// advance past 'if'
	p.advance()

	ifstmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	ifstmt.Then = p.parseBlockStatement()

	if !p.expectPeek(token.RBRACE) {
		p.recover()
		return nil
	}

	if p.peekTokenIs(token.ELSE) || p.peekTokenIs(token.IF) {
		// consume 'else'
		p.advance()
		switch p.curToken.Type {
		case token.IF:
			ifstmt.Else = ast.Stmt(p.parseIfStatement())
		case token.LBRACE:
			ifstmt.Else = ast.Stmt(p.parseBlockStatement())
		}
	}

	return ifstmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStmt {
	exprStmt := &ast.ExpressionStmt{
		Token:      p.curToken,
		Expression: p.parseExpression(LOWEST),
	}

	if !p.expectPeek(token.SEMICOLON) {
		p.recover()
		return nil
	}

	return exprStmt
}

// advances one token and tries to parse an expression based on curToken
func (p *Parser) parseExpression(stickiness int) ast.Expr {
	//   2      +     2
	//   left   op    right
	//   ^
	prefix, ok := p.prefixParselets[p.curToken.Type]

	if !ok {
		p.noParsletError(&p.curToken)
		return nil
	}

	left := prefix()

	for stickiness < p.peekStickiness() {
		//   2      +     2
		//   left   op    right
		//   ^      ^peeking op
		// standing at end of prefix
		// peek next token to see if we can continue
		infix, ok := p.infixParselets[p.peekToken.Type]

		if !ok {
			return left
		}

		p.advance()
		//   2      +     2
		//   left   op    right
		//          ^
		// parse op as infix
		left = infix(left)
	}

	return left
}

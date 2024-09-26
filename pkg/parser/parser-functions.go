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
	case token.LBRACE:
		return p.parseBlockStatement()
	case token.PRINT:
		return p.parsePrintStatement()

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

func (p *Parser) parsePrintStatement() *ast.PrintStmt {
	print := &ast.PrintStmt{
		Token: p.curToken,
	}

	// consunme 'print'
	p.advance()

	print.Expression = p.parseExpression(LOWEST)

	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}

	return print
}

func (p *Parser) parseBlockStatement() *ast.BlockStmt {
	block := &ast.BlockStmt{
		Token:      p.curToken,
		Statements: make([]ast.Stmt, 0),
	}

	// { ... }
	// ^

	// consunme '{'
	p.advance()
	// { ... }
	//   ^

	for !p.curTokenIs(token.RBRACE) && !p.atEnd() {
		stmt := p.parseStatement()
		block.Statements = append(block.Statements, stmt)

		p.advance()
	}

	return block
}

func (p *Parser) parseIfStatement() *ast.IfStmt {
	// if expr { ... } else { ... }
	// ^
	ifstmt := &ast.IfStmt{
		Token: p.curToken,
	}

	// advance past 'if'
	p.advance()

	// if expr { ... } else { ... }
	//    ^
	ifstmt.Condition = p.parseExpression(LOWEST)

	// if expr { ... } else { ... }
	//       ^
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	// if expr { ... } else { ... }
	//         ^
	ifstmt.Then = p.parseBlockStatement()

	// if expr { ... } else { ... }
	//               ^
	if !p.curTokenIs(token.RBRACE) {
		p.recover()
		return nil
	}

	if p.peekTokenIs(token.ELSE) {
		p.advance()
		// if expr { ... } else { ... }
		//                 ^

		p.advance()
		// if expr { ... } else { ... }
		//                      ^
		ifstmt.Else = ast.Stmt(p.parseBlockStatement())

		// if expr { ... } else { ... }
		//                            ^
		if !p.curTokenIs(token.RBRACE) {
			p.recover()
			return nil
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

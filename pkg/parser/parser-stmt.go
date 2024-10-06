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
	case token.IF:
		return p.parseIfStatement()
	case token.LBRACE:
		return p.parseBlockStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.EACH:
		return p.parseEachStatement()
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

	// consume ';' if present. Some expressions dont naturally end with ';'
	p.consume(token.SEMICOLON)

	return letStmt
}

// FIX: does not parse the last element in the list
func (p *Parser) parsePrintStatement() *ast.PrintStmt {
	// print expr1, expr2 ;
	// ^
	print := &ast.PrintStmt{
		Token: p.curToken,
	}

	// print expr1, expr2 ;
	//     ^
	list := p.parseExpressionList(token.SEMICOLON)

	// print expr1, expr2 ;
	//                    ^
	print.Expressions = list

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

func (p *Parser) parseReturnStatement() *ast.ReturnStmt {
	// return ... ;
	// ^
	ret := &ast.ReturnStmt{Token: p.curToken}

	p.advance()
	// return ... ;
	//        ^

	if p.curTokenIs(token.SEMICOLON) || p.curTokenIs(token.RBRACE) {
		// return early with Value as nil
		// return ;
		//        ^
		return ret
	}

	ret.Value = p.parseExpression(LOWEST)

	p.consume(token.SEMICOLON)
	// return ... ;
	//            ^

	return ret
}

// each statments can take different forms
// -- each let a = 0; a < 10; a = a + 1 { ... }
// -- each item in items { ... }

func (p *Parser) parseEachStatement() ast.Stmt {

	// each let a = 0; a < 10; a = a + 1 { ... }
	// ^
	each := &ast.EachStmt{
		Token: p.curToken,
	}

	p.advance()
	// each let a = 0; a < 10; a = a + 1 { ... }
	//      ^
	fmt.Printf("1p.curToken: %s\n", p.curToken.String())
	if !p.curTokenIs(token.SEMICOLON) {
		each.Init = p.parseLetStatment()
	}
	// each let a = 0; a < 10; a = a + 1 { ... }
	//               ^
	p.advance()
	// each let a = 0; a < 10; a = a + 1 { ... }
	//                 ^
	fmt.Printf("2p.curToken: %s\n", p.curToken.String())
	if !p.curTokenIs(token.SEMICOLON) {
		each.Condition = p.parseExpressionStatement().Expression
	}

	// each let a = 0; a < 10; a = a + 1 { ... }
	//                       ^
	p.advance()
	// each let a = 0; a < 10; a = a + 1 { ... }
	//                         ^
	fmt.Printf("3p.curToken: %s\n", p.curToken.String())
	if !p.curTokenIs(token.LBRACE) {
		each.Update = p.parseExpressionStatement().Expression
	}
	// each let a = 0; a < 10; a = a + 1 { ... }
	//                                 ^
	if !p.curTokenIs(token.LBRACE) {
		fmt.Printf("4p.curToken: %s\n", p.curToken.String())
		err := p.expectError(&p.curToken, token.LBRACE)
		p.errors = append(p.errors, err)
		return nil
	}

	fmt.Printf("5p.curToken: %s\n", p.curToken.String())
	body := p.parseBlockStatement()
	fmt.Printf("body.String(): %v\n", body.String())
	each.Body = body
	fmt.Printf("6p.curToken: %s\n", p.curToken.String())

	fmt.Printf("each.String(): %v\n", each.String())
	return each
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStmt {
	exprStmt := &ast.ExpressionStmt{
		Token:      p.curToken,
		Expression: p.parseExpression(LOWEST),
	}

	// consume ';' if present. Some expressions dont naturally end with ';'
	p.consume(token.SEMICOLON)

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

	for !p.peekTokenIs(token.SEMICOLON) && stickiness < p.peekStickiness() {
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

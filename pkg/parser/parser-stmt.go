package parser

import (
	"github.com/fredrikkvalvik/temp-lang/pkg/ast"
	"github.com/fredrikkvalvik/temp-lang/pkg/token"
)

func (p *Parser) parseStatement() ast.Stmt {
	var node ast.Stmt
	switch p.curToken.Type {
	case token.LET:
		node = p.parseLetStatment()
	case token.FUNCTION:
		node = p.parseFunctionStatment()
	case token.IMPORT:
		node = p.parseImportStatement()
	case token.IF:
		node = p.parseIfStatement()
	case token.LBRACE:
		node = p.parseBlockStatement()
	case token.RETURN:
		node = p.parseReturnStatement()
	case token.EACH:
		node = p.parseIteratorStatement()
	case token.PRINT:
		node = p.parsePrintStatement()

	default:
		node = p.parseExpressionStatement()
	}

	if node != nil {
		return node
	} else {
		return nil
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

func (p *Parser) parseImportStatement() *ast.LetStmt {
	return nil
}

// syntactic sugar for declaring a function variable
func (p *Parser) parseFunctionStatment() *ast.LetStmt {
	// fn name ( arg1, arg2 ) { ... }
	// ^
	let := &ast.LetStmt{
		Token: p.curToken,
	}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	// fn name ( arg1, arg2 ) { ... }
	//    ^
	let.Name = &ast.IdentifierExpr{
		Token: p.curToken,
		Value: p.curToken.Lexeme,
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	// fn name ( arg1, arg2 ) { ... }
	//         ^
	fun := &ast.FunctionLiteralExpr{Token: p.curToken}
	fun.Arguments = p.parseFunctionArgs()
	// fn ( arg1, arg2 ) { ... }
	//                 ^

	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	// fn ( arg1, arg2 ) { ... }
	//                   ^

	fun.Body = p.parseBlockStatement()
	// fn ( arg1, arg2 ) { ... }
	//                         ^

	let.Value = fun

	return let
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

func (p *Parser) parseIteratorStatement() *ast.IterStmt {
	// each item : items { ... }
	// ^
	each := &ast.IterStmt{

		Token: p.curToken,
	}

	p.advance()

	// handle case where there is no name or iterable
	// return early
	if p.curTokenIs(token.LBRACE) {
		// each { ... }
		//      ^
		// return with no name or iterable set
		body := p.parseBlockStatement()
		each.Body = body
		return each
	}

	// each items { ... }
	// each item : items { ... }
	//      ^
	first := p.parseExpression(LOWEST)
	// each items { ... }
	// each  item : items { ... }
	//          ^

	if !p.peekTokenIs(token.LBRACE) {
		// set first to iterable and return each with no local var name
		each.Name = first
		p.advance()
		p.advance()
		// each item : items { ... }
		//             ^

		iterable := p.parseExpression(LOWEST)
		// each item : items { ... }
		//                 ^
		each.Iterable = iterable

		p.advance()
		// each item : items { ... }
		//                   ^
	} else {
		each.Iterable = first
		p.advance()
		// each items { ... }
		//            ^
	}

	body := p.parseBlockStatement()
	each.Body = body

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

package parser

import (
	"github.com/fredrikkvalvik/temp-lang/pkg/ast"
	"github.com/fredrikkvalvik/temp-lang/pkg/token"
)

type prefixFn = func() ast.Expr
type infixFn = func(left ast.Expr) ast.Expr

func (p *Parser) registerPrefix(tok token.TokenType, fun prefixFn) {
	p.prefixParselets[tok] = fun
}

func (p *Parser) registerInfix(tok token.TokenType, fun infixFn) {
	p.infixParselets[tok] = fun
}

func (p *Parser) parsePrefix() ast.Expr {
	expr := &ast.UnaryExpr{
		Token:   p.curToken,
		Operand: p.curToken.Type,
	}
	// consume prefix token
	p.advance()

	expr.Right = p.parseExpression(p.peekStickiness())

	return expr
}

func (p *Parser) parseBinary(left ast.Expr) ast.Expr {
	//   2      +     2
	//   left   op    right
	//          ^
	expr := &ast.BinaryExpr{
		Token:   p.curToken,
		Operand: p.curToken.Type,
		Left:    left,
	}

	// we care about how sticky this operator is
	stickiness := p.curStickiness()

	// moving to next operand
	p.advance()
	//   2      +     2
	//   left   op    right
	//                ^

	expr.Right = p.parseExpression(stickiness)

	return expr
}

func (p *Parser) parseParenPrefix() ast.Expr {
	// ( 1 + 2 ) * 3
	// ^
	paren := &ast.ParenExpr{
		Token: p.curToken,
	}

	p.advance()
	// ( 1 + 2 ) * 3
	//   ^

	paren.Expression = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	// ( 1 + 2 ) * 3
	//         ^
	return paren
}

func (p *Parser) parseFunctionLiteral() ast.Expr {
	// fn ( arg1, arg2 ) { ... }
	// ^
	fun := &ast.FunctionLiteralExpr{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	// fn ( arg1, arg2 ) { ... }
	//    ^
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

	return fun
}

func (p *Parser) parseFunctionArgs() []*ast.IdentifierExpr {
	args := []*ast.IdentifierExpr{}

	// fn ( arg1, arg2 ) { ... }
	//    ^

	// handle case with no args
	if p.peekTokenIs(token.RPAREN) {
		p.advance()
		// fn ( arg1, arg2 ) { ... }
		//                 ^
		return args
	}

	p.advance()
	// fn ( arg1, arg2 ) { ... }
	//      ^
	firstArg := &ast.IdentifierExpr{Token: p.curToken, Value: p.curToken.Lexeme}
	args = append(args, firstArg)

	for p.peekTokenIs(token.COMMA) {
		// fn ( arg1, arg2 ) { ... }
		//         ^
		p.advance()
		// fn ( arg1, arg2 ) { ... }
		//          ^
		p.advance()
		// fn ( arg1, arg2 ) { ... }
		//            ^
		arg := &ast.IdentifierExpr{Token: p.curToken, Value: p.curToken.Lexeme}
		args = append(args, arg)
	}

	p.advance()
	// fn ( arg1, arg2 ) { ... }
	//                 ^

	return args

}

func (p *Parser) parseCall(left ast.Expr) ast.Expr {
	fun := &ast.CallExpr{Token: p.curToken}
	fun.Callee = left

	fun.Arguments = p.parseExpressionList(token.RPAREN)

	return fun
}

func (p *Parser) parseIdent() ast.Expr {
	ident := &ast.IdentifierExpr{
		Token: p.curToken,
		Value: p.curToken.Lexeme,
	}

	return ident
}

func (p *Parser) parseNumberLiteral() ast.Expr {
	numberLiteral := &ast.NumberLiteralExpr{
		Token: p.curToken,
		Value: p.curToken.Literal.(float64),
	}

	return numberLiteral
}

func (p *Parser) parseStringLiteral() ast.Expr {
	stringLiteral := &ast.StringLiteralExpr{
		Token: p.curToken,
		Value: p.curToken.Literal.(string),
	}

	return stringLiteral
}

func (p *Parser) parseBooleanLiteral() ast.Expr {
	booleanLiteral := &ast.BooleanLiteralExpr{
		Token: p.curToken,
	}
	if p.curToken.Lexeme == "true" {
		booleanLiteral.Value = true
	} else {
		booleanLiteral.Value = false
	}

	return booleanLiteral
}

func (p *Parser) parseListLiteralExpression() ast.Expr {
	// [ item1, item2 ]
	// ^
	listLiteral := &ast.ListLiteralExpr{
		Token: p.curToken,
	}

	listLiteral.Items = p.parseExpressionList(token.RBRACKET)

	// [ item1, item2 ]
	//                ^
	return listLiteral
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expr {
	list := []ast.Expr{}

	// [ expr1, expr2 ]
	// ^
	if p.peekTokenIs(end) {
		p.advance()
		// [ expr1, expr2 ]
		//                ^
		return list
	}

	p.advance()
	// [ expr1, expr2 ]
	//   ^

	list = append(list, p.parseExpression(LOWEST))
	// [ expr1, expr2 ]
	//       ^

	for p.peekTokenIs(token.COMMA) {
		p.advance()
		// [ expr1, expr2 ]
		//        ^
		p.advance()
		// [ expr1, expr2 ]
		//          ^
		list = append(list, p.parseExpression(LOWEST))
	}

	// [ expr1, expr2 ]
	//              ^
	// error if end of list without seeing `end` token
	if !p.expectPeek(end) {
		return nil
	}
	// [ expr1, expr2 ]
	//                ^

	return list
}

func (p *Parser) parseIndexExpression(left ast.Expr) ast.Expr {
	// list [ expr ]
	//      ^
	expr := &ast.IndexExpr{Token: p.curToken, Left: left}

	p.advance()
	// list [ expr ]
	//        ^
	expr.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}
	// list [ expr ]
	//             ^

	return expr
}

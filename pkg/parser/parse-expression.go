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
		Token: p.curToken,
	}

	expr.Operand = p.parseExpression(p.peekStickiness())

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

	return paren
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

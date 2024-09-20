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

func (p *Parser) precedence() int {
	return 0
}

func (p *Parser) parsePrefix() ast.Expr {
	expr := &ast.PrefixExpr{
		Token: p.curToken,
	}

	expr.Operand = p.parseExpression(p.precedence())

	return expr
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

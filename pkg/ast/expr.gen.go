// THIS FILE IS GENERATED. DO NOT EDIT

package ast

import "github.com/fredrikkvalvik/temp-lang/pkg/token"

type IdentifierExpr struct {
	Token token.Token
	Value string
}

func (n *IdentifierExpr) ExprNode()      {}
func (n *IdentifierExpr) Lexeme() string { return n.Token.Lexeme }
func (n *IdentifierExpr) Literal() any   { return n.Token.Literal }

type NumberLiteralExpr struct {
	Token token.Token
	Value float64
}

func (n *NumberLiteralExpr) ExprNode()      {}
func (n *NumberLiteralExpr) Lexeme() string { return n.Token.Lexeme }
func (n *NumberLiteralExpr) Literal() any   { return n.Token.Literal }

type StringLiteralExpr struct {
	Token token.Token
	Value string
}

func (n *StringLiteralExpr) ExprNode()      {}
func (n *StringLiteralExpr) Lexeme() string { return n.Token.Lexeme }
func (n *StringLiteralExpr) Literal() any   { return n.Token.Literal }

type BooleanLiteralExpr struct {
	Token token.Token
	Value bool
}

func (n *BooleanLiteralExpr) ExprNode()      {}
func (n *BooleanLiteralExpr) Lexeme() string { return n.Token.Lexeme }
func (n *BooleanLiteralExpr) Literal() any   { return n.Token.Literal }

type UnaryExpr struct {
	Token   token.Token
	Operand Expr
}

func (n *UnaryExpr) ExprNode()      {}
func (n *UnaryExpr) Lexeme() string { return n.Token.Lexeme }
func (n *UnaryExpr) Literal() any   { return n.Token.Literal }

type BinaryExpr struct {
	Token   token.Token
	Operand token.TokenType
	Left    Expr
	Right   Expr
}

func (n *BinaryExpr) ExprNode()      {}
func (n *BinaryExpr) Lexeme() string { return n.Token.Lexeme }
func (n *BinaryExpr) Literal() any   { return n.Token.Literal }

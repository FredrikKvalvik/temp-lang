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
	Value float64
	Token token.Token
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

type PrefixExpr struct {
	Token   token.Token
	Operand Expr
}

func (n *PrefixExpr) ExprNode()      {}
func (n *PrefixExpr) Lexeme() string { return n.Token.Lexeme }
func (n *PrefixExpr) Literal() any   { return n.Token.Literal }

type InfixExpr struct {
	Operand Expr
	Token   token.Token
}

func (n *InfixExpr) ExprNode()      {}
func (n *InfixExpr) Lexeme() string { return n.Token.Lexeme }
func (n *InfixExpr) Literal() any   { return n.Token.Literal }

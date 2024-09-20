// THIS FILE IS GENERATED. DO NOT EDIT

package ast

import "github.com/fredrikkvalvik/temp-lang/pkg/token"

type LetStmt struct {
	Token token.Token
	Name  *IdentifierExpr
	Value Expr
}

func (n *LetStmt) StmtNode()      {}
func (n *LetStmt) Lexeme() string { return n.Token.Lexeme }
func (n *LetStmt) Literal() any   { return n.Token.Literal }

type ExpressionStmt struct {
	Token      token.Token
	Expression Expr
}

func (n *ExpressionStmt) StmtNode()      {}
func (n *ExpressionStmt) Lexeme() string { return n.Token.Lexeme }
func (n *ExpressionStmt) Literal() any   { return n.Token.Literal }

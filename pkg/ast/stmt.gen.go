// THIS FILE IS GENERATED. DO NOT EDIT

package ast

import "github.com/fredrikkvalvik/temp-lang/pkg/token"

type LetStmt struct {
	Stmt
	Name  *IdentifierExpr
	Value Expr
	Token token.Token
}

func (n *LetStmt) StmtNode()      {}
func (n *LetStmt) Lexeme() string { return n.Token.Lexeme }
func (n *LetStmt) Literal() any   { return n.Token.Literal }

type IfStmt struct {
	Stmt
	Token     token.Token
	Condition Expr
	Success   Stmt
	Failed    Stmt
}

func (n *IfStmt) StmtNode()      {}
func (n *IfStmt) Lexeme() string { return n.Token.Lexeme }
func (n *IfStmt) Literal() any   { return n.Token.Literal }

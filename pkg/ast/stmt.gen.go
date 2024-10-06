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

type IfStmt struct {
	Token     token.Token
	Condition Expr
	Then      *BlockStmt
	Else      Stmt
}

func (n *IfStmt) StmtNode()      {}
func (n *IfStmt) Lexeme() string { return n.Token.Lexeme }
func (n *IfStmt) Literal() any   { return n.Token.Literal }

type BlockStmt struct {
	Token      token.Token
	Statements []Stmt
}

func (n *BlockStmt) StmtNode()      {}
func (n *BlockStmt) Lexeme() string { return n.Token.Lexeme }
func (n *BlockStmt) Literal() any   { return n.Token.Literal }

type ReturnStmt struct {
	Token token.Token
	Value Expr
}

func (n *ReturnStmt) StmtNode()      {}
func (n *ReturnStmt) Lexeme() string { return n.Token.Lexeme }
func (n *ReturnStmt) Literal() any   { return n.Token.Literal }

type EachStmt struct {
	Token     token.Token
	Init      *LetStmt
	Condition Expr
	Update    Expr
	Body      *BlockStmt
}

func (n *EachStmt) StmtNode()      {}
func (n *EachStmt) Lexeme() string { return n.Token.Lexeme }
func (n *EachStmt) Literal() any   { return n.Token.Literal }

type PrintStmt struct {
	Token       token.Token
	Expressions []Expr
}

func (n *PrintStmt) StmtNode()      {}
func (n *PrintStmt) Lexeme() string { return n.Token.Lexeme }
func (n *PrintStmt) Literal() any   { return n.Token.Literal }

// this is gives us a compile time check to see of all the interafaces has ben properly implemented
func typecheckStmt() {
	_ = Stmt(&LetStmt{})
	_ = Stmt(&ExpressionStmt{})
	_ = Stmt(&IfStmt{})
	_ = Stmt(&BlockStmt{})
	_ = Stmt(&ReturnStmt{})
	_ = Stmt(&EachStmt{})
	_ = Stmt(&PrintStmt{})
}

// THIS FILE IS GENERATED. DO NOT EDIT

package ast

import "github.com/fredrikkvalvik/temp-lang/pkg/token"

type LetStmt struct {
	Token token.Token
	Name  *IdentifierExpr
	Value Expr
}

func (n *LetStmt) StmtNode()              {}
func (n *LetStmt) Lexeme() string         { return n.Token.Lexeme }
func (n *LetStmt) GetToken() *token.Token { return &n.Token }

type ImportStmt struct {
	Token token.Token
	Name  *IdentifierExpr
	Path  string
}

func (n *ImportStmt) StmtNode()              {}
func (n *ImportStmt) Lexeme() string         { return n.Token.Lexeme }
func (n *ImportStmt) GetToken() *token.Token { return &n.Token }

type ExpressionStmt struct {
	Token      token.Token
	Expression Expr
}

func (n *ExpressionStmt) StmtNode()              {}
func (n *ExpressionStmt) Lexeme() string         { return n.Token.Lexeme }
func (n *ExpressionStmt) GetToken() *token.Token { return &n.Token }

type IfStmt struct {
	Token     token.Token
	Condition Expr
	Then      *BlockStmt
	Else      Stmt
}

func (n *IfStmt) StmtNode()              {}
func (n *IfStmt) Lexeme() string         { return n.Token.Lexeme }
func (n *IfStmt) GetToken() *token.Token { return &n.Token }

type BlockStmt struct {
	Token      token.Token
	Statements []Stmt
}

func (n *BlockStmt) StmtNode()              {}
func (n *BlockStmt) Lexeme() string         { return n.Token.Lexeme }
func (n *BlockStmt) GetToken() *token.Token { return &n.Token }

type ReturnStmt struct {
	Token token.Token
	Value Expr
}

func (n *ReturnStmt) StmtNode()              {}
func (n *ReturnStmt) Lexeme() string         { return n.Token.Lexeme }
func (n *ReturnStmt) GetToken() *token.Token { return &n.Token }

type IterStmt struct {
	Token    token.Token
	Name     Expr
	Iterable Expr
	Body     *BlockStmt
}

func (n *IterStmt) StmtNode()              {}
func (n *IterStmt) Lexeme() string         { return n.Token.Lexeme }
func (n *IterStmt) GetToken() *token.Token { return &n.Token }

type WhileStmt struct {
	Token     token.Token
	Condition Expr
	Body      *BlockStmt
}

func (n *WhileStmt) StmtNode()              {}
func (n *WhileStmt) Lexeme() string         { return n.Token.Lexeme }
func (n *WhileStmt) GetToken() *token.Token { return &n.Token }

type PrintStmt struct {
	Token       token.Token
	Expressions []Expr
}

func (n *PrintStmt) StmtNode()              {}
func (n *PrintStmt) Lexeme() string         { return n.Token.Lexeme }
func (n *PrintStmt) GetToken() *token.Token { return &n.Token }

// this is gives us a compile time check to see of all the interafaces has ben properly implemented
func _() {
	_ = Stmt(&LetStmt{})
	_ = Stmt(&ImportStmt{})
	_ = Stmt(&ExpressionStmt{})
	_ = Stmt(&IfStmt{})
	_ = Stmt(&BlockStmt{})
	_ = Stmt(&ReturnStmt{})
	_ = Stmt(&IterStmt{})
	_ = Stmt(&WhileStmt{})
	_ = Stmt(&PrintStmt{})
}

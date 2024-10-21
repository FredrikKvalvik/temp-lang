// THIS FILE IS GENERATED. DO NOT EDIT

package ast

import "github.com/fredrikkvalvik/temp-lang/pkg/token"

type IdentifierExpr struct {
	Token           token.Token
	Value           string
	ResolutionDepth int
}

func (n *IdentifierExpr) ExprNode()              {}
func (n *IdentifierExpr) Lexeme() string         { return n.Token.Lexeme }
func (n *IdentifierExpr) GetToken() *token.Token { return &n.Token }

type NumberLiteralExpr struct {
	Token token.Token
	Value float64
}

func (n *NumberLiteralExpr) ExprNode()              {}
func (n *NumberLiteralExpr) Lexeme() string         { return n.Token.Lexeme }
func (n *NumberLiteralExpr) GetToken() *token.Token { return &n.Token }

type StringLiteralExpr struct {
	Token token.Token
	Value string
}

func (n *StringLiteralExpr) ExprNode()              {}
func (n *StringLiteralExpr) Lexeme() string         { return n.Token.Lexeme }
func (n *StringLiteralExpr) GetToken() *token.Token { return &n.Token }

type BooleanLiteralExpr struct {
	Token token.Token
	Value bool
}

func (n *BooleanLiteralExpr) ExprNode()              {}
func (n *BooleanLiteralExpr) Lexeme() string         { return n.Token.Lexeme }
func (n *BooleanLiteralExpr) GetToken() *token.Token { return &n.Token }

type UnaryExpr struct {
	Token   token.Token
	Operand token.TokenType
	Right   Expr
}

func (n *UnaryExpr) ExprNode()              {}
func (n *UnaryExpr) Lexeme() string         { return n.Token.Lexeme }
func (n *UnaryExpr) GetToken() *token.Token { return &n.Token }

type BinaryExpr struct {
	Token   token.Token
	Operand token.TokenType
	Left    Expr
	Right   Expr
}

func (n *BinaryExpr) ExprNode()              {}
func (n *BinaryExpr) Lexeme() string         { return n.Token.Lexeme }
func (n *BinaryExpr) GetToken() *token.Token { return &n.Token }

type LogicalExpr struct {
	Token   token.Token
	Operand token.TokenType
	Left    Expr
	Right   Expr
}

func (n *LogicalExpr) ExprNode()              {}
func (n *LogicalExpr) Lexeme() string         { return n.Token.Lexeme }
func (n *LogicalExpr) GetToken() *token.Token { return &n.Token }

type AssignExpr struct {
	Token    token.Token
	Operand  token.TokenType
	Assignee Expr
	Value    Expr
}

func (n *AssignExpr) ExprNode()              {}
func (n *AssignExpr) Lexeme() string         { return n.Token.Lexeme }
func (n *AssignExpr) GetToken() *token.Token { return &n.Token }

type ParenExpr struct {
	Token      token.Token
	Expression Expr
}

func (n *ParenExpr) ExprNode()              {}
func (n *ParenExpr) Lexeme() string         { return n.Token.Lexeme }
func (n *ParenExpr) GetToken() *token.Token { return &n.Token }

type FunctionLiteralExpr struct {
	Token     token.Token
	Arguments []*IdentifierExpr
	Body      *BlockStmt
}

func (n *FunctionLiteralExpr) ExprNode()              {}
func (n *FunctionLiteralExpr) Lexeme() string         { return n.Token.Lexeme }
func (n *FunctionLiteralExpr) GetToken() *token.Token { return &n.Token }

type CallExpr struct {
	Token     token.Token
	Callee    Expr
	Arguments []Expr
}

func (n *CallExpr) ExprNode()              {}
func (n *CallExpr) Lexeme() string         { return n.Token.Lexeme }
func (n *CallExpr) GetToken() *token.Token { return &n.Token }

type ListLiteralExpr struct {
	Token token.Token
	Items []Expr
}

func (n *ListLiteralExpr) ExprNode()              {}
func (n *ListLiteralExpr) Lexeme() string         { return n.Token.Lexeme }
func (n *ListLiteralExpr) GetToken() *token.Token { return &n.Token }

type MapLiteralExpr struct {
	Token     token.Token
	KeyValues map[Expr]Expr
}

func (n *MapLiteralExpr) ExprNode()              {}
func (n *MapLiteralExpr) Lexeme() string         { return n.Token.Lexeme }
func (n *MapLiteralExpr) GetToken() *token.Token { return &n.Token }

type IndexExpr struct {
	Token token.Token
	Left  Expr
	Index Expr
}

func (n *IndexExpr) ExprNode()              {}
func (n *IndexExpr) Lexeme() string         { return n.Token.Lexeme }
func (n *IndexExpr) GetToken() *token.Token { return &n.Token }

// this is gives us a compile time check to see of all the interafaces has ben properly implemented
func _() {
	_ = Expr(&IdentifierExpr{})
	_ = Expr(&NumberLiteralExpr{})
	_ = Expr(&StringLiteralExpr{})
	_ = Expr(&BooleanLiteralExpr{})
	_ = Expr(&UnaryExpr{})
	_ = Expr(&BinaryExpr{})
	_ = Expr(&LogicalExpr{})
	_ = Expr(&AssignExpr{})
	_ = Expr(&ParenExpr{})
	_ = Expr(&FunctionLiteralExpr{})
	_ = Expr(&CallExpr{})
	_ = Expr(&ListLiteralExpr{})
	_ = Expr(&MapLiteralExpr{})
	_ = Expr(&IndexExpr{})
}

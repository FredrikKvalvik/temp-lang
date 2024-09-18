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

//go:generate go run generate-ast-nodes.go

/*
the ast package hold the ast nodes for `temp-lang` for the ast that will be
generated by a parser, and evaluated as a program
*/
package ast

import (
	"fmt"
	"strings"

	"github.com/fredrikkvalvik/temp-lang/pkg/token"
)

type Node interface {
	Literal() any
	Lexeme() string
	String() string
	GetToken() *token.Token
}

type Stmt interface {
	Node
	StmtNode() // Dummy method to make go treat Expr and Stmt differently
}

type Expr interface {
	Node
	ExprNode() // Dummy method to make go treat Expr and Stmt differently
}

type Program struct {
	Statements []Stmt
}

func (p *Program) Literal() any   { return "PROGRAM" }
func (p *Program) Lexeme() string { return "PROGRAM" }
func (p *Program) GetToken() *token.Token {
	return &token.Token{Type: token.ILLEGAL}
}
func (p *Program) String() string {
	var str strings.Builder

	for _, stmt := range p.Statements {
		str.WriteString(stmt.String())
	}

	return str.String()
}

// / =========
func (l *LetStmt) String() string {
	var s strings.Builder

	// TODO: update String when let is is fully implemented
	fmt.Fprintf(&s, "let %s = %s;\n", l.Name.String(), l.Value.String())

	return s.String()
}

func (e *ExpressionStmt) String() string {
	var s strings.Builder

	fmt.Fprintf(&s, "%s\n", e.Expression.String())

	return s.String()
}

func (e *PrintStmt) String() string {
	var s strings.Builder

	for _, expr := range e.Expressions {
		fmt.Fprintf(&s, "print %s ", expr.String())
	}

	return s.String()
}

func (i *IfStmt) String() string {
	var s strings.Builder

	fmt.Fprintf(&s, "if %s ", i.Condition.String())
	fmt.Fprint(&s, i.Then.String())

	if i.Else != nil {
		fmt.Fprintf(&s, " else ")
		fmt.Fprint(&s, i.Else.String())

	}
	return s.String()
}

func (b *BlockStmt) String() string {
	var s strings.Builder

	s.WriteString("{\n")
	for _, stmt := range b.Statements {
		s.WriteString(stmt.String())
	}
	s.WriteString("}")

	return s.String()
}

func (r *ReturnStmt) String() string {
	str := "return"

	if r.Value != nil {
		str += " " + r.Token.String()
	}

	return str
}

func (r *EachStmt) String() string {
	var str strings.Builder

	init := ""
	if r.Init != nil {
		init = r.Init.String()
	}
	condition := ""
	if r.Condition != nil {
		condition = r.Condition.String()
	}
	update := ""
	if r.Condition != nil {
		update = r.Update.String()
	}

	fmt.Fprintf(
		&str,
		"each %s; %s; %s %s",
		init,
		condition,
		update,
		r.Body.String(),
	)

	return str.String()
}

func (s *IterStmt) String() string {
	var str strings.Builder

	fmt.Fprint(&str, "each")

	switch {
	case s.Iterable == nil && s.Name == nil:
		fmt.Fprintf(&str, " ")
	case s.Iterable != nil && s.Name == nil:
		fmt.Fprintf(&str, " %s ", s.Iterable.String())
	case s.Iterable != nil && s.Name != nil:
		fmt.Fprintf(&str, " %s : %s ", s.Name.String(), s.Iterable.String())
	}

	fmt.Fprint(&str, s.Body.String())

	return str.String()
}

func (i *IdentifierExpr) String() string {
	var str strings.Builder

	fmt.Fprintf(&str, "%s", i.Value)

	return str.String()
}

func (p *UnaryExpr) String() string {
	var s strings.Builder

	fmt.Fprintf(&s, "(%s)", p.Lexeme()+p.Right.String())

	return s.String()
}

func (b *BinaryExpr) String() string {
	var s strings.Builder

	fmt.Fprintf(&s, "(%s %s %s)", b.Left.String(), b.Lexeme(), b.Right.String())

	return s.String()

}

func (p *ParenExpr) String() string {
	var s strings.Builder

	fmt.Fprint(&s, p.Expression.String())

	return s.String()

}

func (p *FunctionLiteralExpr) String() string {
	var s strings.Builder

	s.WriteString("fn(")
	for idx, arg := range p.Arguments {
		s.WriteString(arg.Value)

		if len(p.Arguments) != idx+1 {
			s.WriteString(", ")
		}
	}
	fmt.Fprintf(&s, ") %s", p.Body.String())

	return s.String()
}

func (n *CallExpr) String() string {
	var s strings.Builder

	s.WriteString(fmt.Sprintf("fn %s(", n.Callee.Lexeme()))
	for i, arg := range n.Arguments {
		s.WriteString(arg.String())

		if i != len(n.Arguments)-1 {
			s.WriteString(", ")
		}
	}
	s.WriteString(")")

	return s.String()
}

func (n *IndexExpr) String() string {
	return "TODO INDEX STRING"
}

func (n *ListLiteralExpr) String() string {
	return "LIST LIT TODO"
}

func (n *NumberLiteralExpr) String() string {
	return n.Lexeme()
}

func (s *StringLiteralExpr) String() string {
	return s.Lexeme()
}

func (s *BooleanLiteralExpr) String() string {
	return s.Lexeme()
}

// HELPERS

// func indent(input string, indent int) string {
// 	lines := strings.Split(input, "\n")

// 	for i, l := range lines {
// 		newStr := strings.Repeat(" ", indent) + l
// 		lines[i] = newStr
// 	}

// 	return strings.Join(lines, "\n")
// }

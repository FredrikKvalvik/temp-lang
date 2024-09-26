//go:generate go run generate-ast-nodes.go

/*
the ast package hold the ast nodes for `temp-lang` for the ast that will be
generated by a parser, and evaluated as a program
*/
package ast

import (
	"fmt"
	"strings"
)

type Node interface {
	Literal() any
	Lexeme() string
	String() string
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
	s.WriteString(")")

	return s.String()

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

//go:build ignore

package main

import (
	"fmt"
	"os"
	"strings"
)

const stmt = "Stmt"
const expr = "Expr"

const stmtMethod = "func (%s *%s) StmtNode() {}"
const exprMethod = "func (%s *%s) ExprNode() {}"

const packageName = "ast"
const tokenPkg = "github.com/fredrikkvalvik/temp-lang/pkg/token"

type template struct {
	name  string
	props map[string]string
}

var stmts = []template{
	{
		name: "Let",
		props: map[string]string{
			"Token": "token.Token",
			"Name":  "*Identifier" + expr,
			"Value": expr,
		},
	},
	{
		name: "Expression",
		props: map[string]string{
			"Token":      "token.Token",
			"Expression": expr,
		},
	},
}

var exprs = []template{
	{
		name: "Identifier",
		props: map[string]string{
			"Token": "token.Token",
			"Value": "string",
		},
	},
	{
		name: "NumberLiteral",
		props: map[string]string{
			"Token": "token.Token",
			"Value": "float64",
		},
	},
	{
		name: "StringLiteral",
		props: map[string]string{
			"Token": "token.Token",
			"Value": "string",
		},
	},
	{
		name: "BooleanLiteral",
		props: map[string]string{
			"Token": "token.Token",
			"Value": "bool",
		},
	},
	{
		name: "Prefix",
		props: map[string]string{
			"Token":   "token.Token",
			"Operand": expr,
		},
	},
	{
		name: "Infix",
		props: map[string]string{
			"Token":   "token.Token",
			"Operand": expr,
		},
	},
}

// This will generate a file for statements and expressions
// the only unique part of the structs are the fields and the String method
func main() {
	statementsFile := generateNodes(stmt, stmtMethod, stmts)
	expressionFile := generateNodes(expr, exprMethod, exprs)

	os.WriteFile("stmt.gen.go", []byte(statementsFile), 0646)
	os.WriteFile("expr.gen.go", []byte(expressionFile), 0646)
}

func generateNodes(interfaceName, interfaceMethod string, tmpl []template) string {
	var f strings.Builder

	f.WriteString("// THIS FILE IS GENERATED. DO NOT EDIT\n\n")
	f.WriteString(fmt.Sprintf("package %s\n\n", packageName))
	f.WriteString(fmt.Sprintf(`import "%s"`+"\n\n", tokenPkg))

	for _, s := range tmpl {
		name := s.name + interfaceName

		f.WriteString(fmt.Sprintf("type %s struct {\n", name))

		// f.WriteString(fmt.Sprintf("\t%s\n", interfaceName))
		for key, value := range s.props {
			f.WriteString(fmt.Sprintf("\t%s %s\n", key, value))
		}

		f.WriteString("}\n")
		f.WriteString(fmt.Sprintf(interfaceMethod, "n", name) + "\n")
		f.WriteString(fmt.Sprintf("func (n *%s) Lexeme() string { return n.Token.Lexeme }\n", name))
		f.WriteString(fmt.Sprintf("func (n *%s) Literal() any { return n.Token.Literal }\n", name))

		// create space for next struct
		f.WriteString("\n")
	}

	return f.String()
}

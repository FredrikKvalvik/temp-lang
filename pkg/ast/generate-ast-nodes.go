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

type keyVal struct {
	key   string
	value string
}
type template struct {
	name  string
	props []keyVal
}

var stmts = []template{
	{
		name: "Let",
		props: []keyVal{
			{"Token", "token.Token"},
			{"Name", "*Identifier" + expr},
			{"Value", expr},
		},
	},
	{
		name: "Expression",
		props: []keyVal{
			{"Token", "token.Token"},
			{"Expression", expr},
		},
	},
	{
		name: "If",
		props: []keyVal{
			{"Token", "token.Token"},
			{"Condition", expr},
			{"Then", "*Block" + stmt},
			{"Else", stmt},
		},
	},
	{
		name: "Block",
		props: []keyVal{
			{"Token", "token.Token"},
			{"Statements", "[]" + stmt},
		},
	},
	{
		name: "Return",
		props: []keyVal{
			{"Token", "token.Token"},
			{"Value", expr},
		},
	},
	{
		name: "Print",
		props: []keyVal{
			{"Token", "token.Token"},
			{"Expressions", "[]" + expr},
		},
	},
}

var exprs = []template{
	{
		name: "Identifier",
		props: []keyVal{
			{"Token", "token.Token"},
			{"Value", "string"},
		},
	},
	{
		name: "NumberLiteral",
		props: []keyVal{
			{"Token", "token.Token"},
			{"Value", "float64"},
		},
	},
	{
		name: "StringLiteral",
		props: []keyVal{
			{"Token", "token.Token"},
			{"Value", "string"},
		},
	},
	{
		name: "BooleanLiteral",
		props: []keyVal{
			{"Token", "token.Token"},
			{"Value", "bool"},
		},
	},
	{
		name: "Unary",
		props: []keyVal{
			{"Token", "token.Token"},
			{"Operand", "token.TokenType"},
			{"Right", expr},
		},
	},
	{
		name: "Binary",
		props: []keyVal{
			{"Token", "token.Token"},
			{"Operand", "token.TokenType"},
			{"Left", expr},
			{"Right", expr},
		},
	},
	{
		name: "Paren",
		props: []keyVal{
			{"Token", "token.Token"},
			{"Expression", expr},
		},
	},
	{
		name: "FunctionLiteral",
		props: []keyVal{
			{"Token", "token.Token"},
			{"Arguments", "[]*Identifier" + expr},
			{"Body", "*Block" + stmt},
		},
	},
	{
		name: "Call",
		props: []keyVal{
			{"Token", "token.Token"},
			{"Callee", expr},
			{"Arguments", "[]" + expr},
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
		for _, kv := range s.props {
			f.WriteString(fmt.Sprintf("\t%s %s\n", kv.key, kv.value))
		}

		f.WriteString("}\n")
		f.WriteString(fmt.Sprintf(interfaceMethod, "n", name) + "\n")
		f.WriteString(fmt.Sprintf("func (n *%s) Lexeme() string { return n.Token.Lexeme }\n", name))
		f.WriteString(fmt.Sprintf("func (n *%s) Literal() any { return n.Token.Literal }\n", name))

		// create space for next struct
		f.WriteString("\n")
	}

	fmt.Fprint(&f, "// this is gives us a compile time check to see of all the interafaces has ben properly implemented\n")
	fmt.Fprintf(&f, "func typecheck%s() {\n", interfaceName)
	for _, s := range tmpl {
		name := s.name + interfaceName

		fmt.Fprintf(&f, "_ = %s(&%s{})\n", interfaceName, name)
	}
	fmt.Fprint(&f, "}")

	return f.String()
}

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
			{"Name", "*Identifier" + expr},
			{"Value", expr},
		},
	},
	{
		name: "Expression",
		props: []keyVal{
			{"Expression", expr},
		},
	},
	{
		name: "If",
		props: []keyVal{
			{"Condition", expr},
			{"Then", "*Block" + stmt},
			{"Else", stmt},
		},
	},
	{
		name: "Block",
		props: []keyVal{
			{"Statements", "[]" + stmt},
		},
	},
	{
		name: "Return",
		props: []keyVal{
			{"Value", expr},
		},
	},
	{
		name: "Iter",
		props: []keyVal{
			{"Name", expr},
			{"Iterable", expr},
			{"Body", "*Block" + stmt},
		},
	},
	{
		name: "Print",
		props: []keyVal{
			{"Expressions", "[]" + expr},
		},
	},
}

var exprs = []template{
	{
		name: "Identifier",
		props: []keyVal{
			{"Value", "string"},
		},
	},
	{
		name: "NumberLiteral",
		props: []keyVal{
			{"Value", "float64"},
		},
	},
	{
		name: "StringLiteral",
		props: []keyVal{
			{"Value", "string"},
		},
	},
	{
		name: "BooleanLiteral",
		props: []keyVal{
			{"Value", "bool"},
		},
	},
	{
		name: "Unary",
		props: []keyVal{
			{"Operand", "token.TokenType"},
			{"Right", expr},
		},
	},
	{
		name: "Binary",
		props: []keyVal{
			{"Operand", "token.TokenType"},
			{"Left", expr},
			{"Right", expr},
		},
	},
	{
		name: "Paren",
		props: []keyVal{
			{"Expression", expr},
		},
	},
	{
		name: "FunctionLiteral",
		props: []keyVal{
			{"Arguments", "[]*Identifier" + expr},
			{"Body", "*Block" + stmt},
		},
	},
	{
		name: "Call",
		props: []keyVal{
			{"Callee", expr},
			{"Arguments", "[]" + expr},
		},
	},
	{
		name: "ListLiteral",
		props: []keyVal{
			{"Items", "[]" + expr},
		},
	},
	{
		name: "MapLiteral",
		props: []keyVal{
			{"KeyValues", fmt.Sprintf("map[%s]%s", expr, expr)},
		},
	},
	{
		name: "Index",
		props: []keyVal{
			{"Left", expr},
			{"Index", expr},
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

		// add token to all ast nodes
		kvs := append([]keyVal{{"Token", "token.Token"}}, s.props...)
		for _, kv := range kvs {
			f.WriteString(fmt.Sprintf("\t%s %s\n", kv.key, kv.value))
		}

		f.WriteString("}\n")
		f.WriteString(fmt.Sprintf(interfaceMethod, "n", name) + "\n")
		f.WriteString(fmt.Sprintf("func (n *%s) Lexeme() string { return n.Token.Lexeme }\n", name))
		f.WriteString(fmt.Sprintf("func (n *%s) Literal() any { return n.Token.Literal }\n", name))
		f.WriteString(fmt.Sprintf("func (n *%s) GetToken() *token.Token { return &n.Token }\n", name))

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

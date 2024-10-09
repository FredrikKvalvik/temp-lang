//go:build ignore

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fredrikkvalvik/temp-lang/pkg/object"
)

const packageName = "object"
const astPkg = "github.com/fredrikkvalvik/temp-lang/pkg/ast"
const tokenPkg = "github.com/fredrikkvalvik/temp-lang/pkg/token"

type keyVal struct {
	key   string
	value string
}
type template struct {
	name  string
	typ   object.ObjectType
	props []keyVal
}

var objects = []template{
	{
		name: "Boolean",
		typ:  object.BOOL_OBJ,
		props: []keyVal{
			{"Value", "bool"},
		},
	},
	{
		name:  "Nil",
		props: []keyVal{},
		typ:   object.NIL_OBJ,
	},
	{
		name: "Number",
		typ:  object.NUMBER_OBJ,
		props: []keyVal{
			{"Value", "float64"},
		},
	},
	{
		name: "String",
		typ:  object.STRING_OBJ,
		props: []keyVal{
			{"Value", "string"},
		},
	},
	{
		name: "FnLiteral",
		typ:  object.FUNCTION_LITERAL_OBJ,
		props: []keyVal{
			{"Parameters", "[]*ast.IdentifierExpr"},
			{"Body", "*ast.BlockStmt"},
			{"Env", "*Environment"},
		},
	},
	{
		name: "Return",
		typ:  object.RETURN_OBJ,
		props: []keyVal{
			{"Value", "Object"},
		},
	},
	{
		name: "List",
		typ:  object.LIST_OBJ,
		props: []keyVal{
			{"Values", "[]Object"},
		},
	},
	{
		name: "Error",
		typ:  object.ERROR_OBJ,
		props: []keyVal{
			{"Error", "error"},
			{"Token", "token.Token"},
		},
	},
}

// This will generate a file for statements and expressions
// the only unique part of the structs are the fields and the String method
func main() {
	objectFile := generateObjects(objects)

	os.WriteFile("objects.gen.go", []byte(objectFile), 0646)
}

func generateObjects(tmpl []template) string {
	var f strings.Builder

	f.WriteString("// THIS FILE IS GENERATED. DO NOT EDIT\n\n")
	f.WriteString(fmt.Sprintf("package %s\n\n", packageName))
	f.WriteString(fmt.Sprintf(`import "%s"`+"\n\n", astPkg))
	f.WriteString(fmt.Sprintf(`import "%s"`+"\n\n", tokenPkg))

	for _, s := range tmpl {
		name := s.name + "Obj"

		f.WriteString(fmt.Sprintf("type %s struct {\n", name))

		// f.WriteString(fmt.Sprintf("\t%s\n", interfaceName))
		for _, kv := range s.props {
			f.WriteString(fmt.Sprintf("\t%s %s\n", kv.key, kv.value))
		}

		f.WriteString("}\n")
		f.WriteString(fmt.Sprintf("func (n *%s) Type() ObjectType { return %s }\n", name, s.typ))

		// create space for next struct
		f.WriteString("\n")
	}

	fmt.Fprint(&f, "// this is gives us a compile time check to see of all the interafaces has been properly implemented\n")
	fmt.Fprintf(&f, "func typecheck() {\n")
	for _, s := range tmpl {
		name := s.name + "Obj"

		fmt.Fprintf(&f, "_ = Object(& %s {})\n", name)
	}
	fmt.Fprint(&f, "}")

	return f.String()
}

//go:build ignore

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fredrikkvalvik/temp-lang/pkg/object"
)

const packageName = "object"

type template struct {
	name  string
	typ   object.ObjectType
	props map[string]string
}

var objects = []template{
	{
		name: "Boolean",
		typ:  object.BOOL_OBJ,
		props: map[string]string{
			"Value": "bool",
		},
	},
	{
		name:  "Nil",
		typ:   object.NIL_OBJ,
		props: map[string]string{},
	},
	{
		name: "Number",
		typ:  object.NUMBER_OBJ,
		props: map[string]string{
			"Value": "float64",
		},
	},
	{
		name: "String",
		typ:  object.STRING_OBJ,
		props: map[string]string{
			"Value": "string",
		},
	},
	{
		name: "Error",
		typ:  object.ERROR_OBJ,
		props: map[string]string{
			"Message": "string",
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

	for _, s := range tmpl {
		name := s.name

		f.WriteString(fmt.Sprintf("type %s struct {\n", name))

		// f.WriteString(fmt.Sprintf("\t%s\n", interfaceName))
		for key, value := range s.props {
			f.WriteString(fmt.Sprintf("\t%s %s\n", key, value))
		}

		f.WriteString("}\n")
		f.WriteString(fmt.Sprintf("func (n *%s) Type() ObjectType { return %s }\n", name, s.typ))

		// create space for next struct
		f.WriteString("\n")
	}

	fmt.Fprint(&f, "// this is gives us a compile time check to see of all the interafaces has been properly implemented\n")
	fmt.Fprintf(&f, "func typecheck() {\n")
	for _, s := range tmpl {
		name := s.name

		fmt.Fprintf(&f, "_ = Object(& %s {})\n", name)
	}
	fmt.Fprint(&f, "}")

	return f.String()
}
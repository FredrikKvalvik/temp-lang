//go:generate go run golang.org/x/tools/cmd/stringer -type=ObjectType
//go:generate go run generate-objects.go

package object

import (
	"fmt"
	"strings"
)

// object represents runtime values.
// Object can be any value thats valid in the program
type Object interface {
	Type() ObjectType
	Inspect() string
}

type BuiltinFn func(args ...Object) Object

type MapPairs = map[HashKey]KeyValuePair

type ObjectType int

const (
	_ ObjectType = iota

	OBJ_BOOL             // representes true and false
	OBJ_NIL              // sentinel value for "no value"
	OBJ_NUMBER           // number object is any float64-representable number
	OBJ_STRING           // represents a string value
	OBJ_FUNCTION_LITERAL // represents a function object
	OBJ_RETURN           // internal type for propagating return values
	OBJ_LIST             // collection if objects in an ordered list
	OBJ_MAP              // Map is a datatype for storing key-value pairs
	OBJ_BUILTIN          // Builtin function
	OBJ_ITERATOR         // a wrapper for returning iterators from builtin functions
	OBJ_MODULE           // Module is an object that holds the references to a unit of code that has been imported by a caller
	OBJ_ERROR            // runtime error
)

type ModuleType int

const (
	_ ModuleType = iota
	NATIVE_MODULE
	FILE_MODULE
)

func (n *NilObj) Inspect() string     { return "nil" }
func (b *BooleanObj) Inspect() string { return fmt.Sprintf("%v", b.Value) }
func (b *StringObj) Inspect() string  { return fmt.Sprintf(`"%s"`, b.Value) }
func (b *NumberObj) Inspect() string  { return fmt.Sprintf("%v", b.Value) }
func (b *FnLiteralObj) Inspect() string {
	var str strings.Builder

	str.WriteString("fn(")
	for idx, param := range b.Parameters {
		if idx != 0 {
			str.WriteString(", ")
		}

		str.WriteString(param.String())
	}

	str.WriteString(") ")
	str.WriteString(b.Body.String())

	return str.String()
}
func (b *ReturnObj) Inspect() string { return fmt.Sprintf("return[%s]", b.Value.Inspect()) }
func (b *ListObj) Inspect() string {
	var str strings.Builder

	fmt.Fprint(&str, "[")

	for i, value := range b.Values {
		if i != len(b.Values)-1 {
			fmt.Fprintf(&str, "%s, ", value.Inspect())
		} else {
			fmt.Fprint(&str, value.Inspect())
		}
	}

	fmt.Fprint(&str, "]")

	return str.String()
}

func (b *MapObj) Inspect() string {
	var str strings.Builder

	fmt.Fprint(&str, "{\n")

	for _, kv := range b.Pairs {
		fmt.Fprintf(&str, "  %s: %s,\n", kv.Key.Inspect(), kv.Value.Inspect())
	}

	fmt.Fprint(&str, "}")

	return str.String()
}

func (b *IteratorObj) Inspect() string { return fmt.Sprintf("[Iterator %s]", b.Iterator.Type()) }

func (b *BuiltinObj) Inspect() string { return fmt.Sprintf("[builtin %s]", b.Name) }

func (b *ModuleObj) Inspect() string { return fmt.Sprintf("[Module %s]", b.Name) }

func (b *ErrorObj) Inspect() string {
	if b.Token != nil {
		line, col := b.Token.Pos.Position()
		return fmt.Sprintf("[%d:%d] %s", line, col, b.Error.Error())
	} else {
		return b.Error.Error()
	}
}

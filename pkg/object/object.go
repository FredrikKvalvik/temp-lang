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
type ObjectType int

type MapPairs = map[HashKey]KeyValuePair

const (
	_                    ObjectType = iota
	BOOL_OBJ                        // representes true and false
	NIL_OBJ                         // sentinel value for "no value"
	NUMBER_OBJ                      // number object is any float64-representable number
	STRING_OBJ                      // represents a string value
	FUNCTION_LITERAL_OBJ            // represents a function object
	RETURN_OBJ                      // internal type for propagating return values
	LIST_OBJ                        // collection if objects in an ordered list
	MAP_OBJ                         // Map is a datatype for storing key-value pairs
	BUILTIN_OBJ                     // Builtin function
	ERROR_OBJ                       // runtime error
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
func (b *BuiltinObj) Inspect() string { return fmt.Sprintf("[builtin %s]", b.Name) }

func (b *ErrorObj) Inspect() string { return b.Error.Error() }

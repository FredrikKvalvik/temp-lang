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

type ObjectType int

const (
	BOOL_OBJ ObjectType = iota
	NIL_OBJ
	NUMBER_OBJ
	STRING_OBJ
	FUNCTION_LITERAL_OBJ
	RETURN_OBJ
	LIST_OBJ
	ERROR_OBJ
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
func (b *ErrorObj) Inspect() string { return b.Error.Error() }

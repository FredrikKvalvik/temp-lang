//go:generate go run golang.org/x/tools/cmd/stringer -type=ObjectType
//go:generate go run generate-objects.go

package object

import "fmt"

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
	ERROR_OBJ
)

func (n *Nil) Inspect() string       { return "nil" }
func (b *Boolean) Inspect() string   { return fmt.Sprintf("%v", b.Value) }
func (b *String) Inspect() string    { return fmt.Sprintf(`"%s"`, b.Value) }
func (b *Number) Inspect() string    { return fmt.Sprintf("%v", b.Value) }
func (b *FnLiteral) Inspect() string { return fmt.Sprint("FnLiteral inspect not implemented") }
func (b *Error) Inspect() string     { return b.Message }

//go:generate go run golang.org/x/tools/cmd/stringer -type=ObjectType
//go:generate go run generate-objects.go

package object

// object represents runtime values.
// Object can be any value thats valid in the program
type Object interface {
	Type() ObjectType
	Inspect() string
}

type ObjectType int

const (
	BOOL_OBJ ObjectType = iota
	NULL_OBJ
	NUMBER_OBJ
	STRING_OBJ
)

func (n *Null) Inspect() string    { return "NOT IMPLEMENTED" }
func (b *Boolean) Inspect() string { return "NOT IMPLEMENTED" }
func (b *String) Inspect() string  { return "NOT IMPLEMENTED" }
func (b *Number) Inspect() string  { return "NOT IMPLEMENTED" }

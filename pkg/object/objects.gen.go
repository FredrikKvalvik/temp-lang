// THIS FILE IS GENERATED. DO NOT EDIT

package object

type Boolean struct {
	Value bool
}

func (n *Boolean) Type() ObjectType { return BOOL_OBJ }

type Null struct {
}

func (n *Null) Type() ObjectType { return NULL_OBJ }

type Number struct {
	Value string
}

func (n *Number) Type() ObjectType { return NUMBER_OBJ }

type String struct {
	Value string
}

func (n *String) Type() ObjectType { return STRING_OBJ }

// this is gives us a compile time check to see of all the interafaces has ben properly implemented
func typecheck() {
	_ = Object(&Boolean{})
	_ = Object(&Null{})
	_ = Object(&Number{})
	_ = Object(&String{})
}

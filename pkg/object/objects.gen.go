// THIS FILE IS GENERATED. DO NOT EDIT

package object

type Boolean struct {
	Value bool
}

func (n *Boolean) Type() ObjectType { return BOOL_OBJ }

type Nil struct {
}

func (n *Nil) Type() ObjectType { return NIL_OBJ }

type Number struct {
	Value float64
}

func (n *Number) Type() ObjectType { return NUMBER_OBJ }

type String struct {
	Value string
}

func (n *String) Type() ObjectType { return STRING_OBJ }

type Error struct {
	Message string
}

func (n *Error) Type() ObjectType { return ERROR_OBJ }

// this is gives us a compile time check to see of all the interafaces has been properly implemented
func typecheck() {
	_ = Object(&Boolean{})
	_ = Object(&Nil{})
	_ = Object(&Number{})
	_ = Object(&String{})
	_ = Object(&Error{})
}

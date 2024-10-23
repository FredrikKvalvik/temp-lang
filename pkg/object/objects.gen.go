// THIS FILE IS GENERATED. DO NOT EDIT

package object

import "github.com/fredrikkvalvik/temp-lang/pkg/ast"

import "github.com/fredrikkvalvik/temp-lang/pkg/token"

type BooleanObj struct {
	Value bool
}

func (n *BooleanObj) Type() ObjectType { return BOOL_OBJ }

type NilObj struct {
}

func (n *NilObj) Type() ObjectType { return NIL_OBJ }

type NumberObj struct {
	Value float64
}

func (n *NumberObj) Type() ObjectType { return NUMBER_OBJ }

type StringObj struct {
	Value string
}

func (n *StringObj) Type() ObjectType { return STRING_OBJ }

type FnLiteralObj struct {
	Parameters []*ast.IdentifierExpr
	Body       *ast.BlockStmt
	Env        *Environment
}

func (n *FnLiteralObj) Type() ObjectType { return FUNCTION_LITERAL_OBJ }

type ReturnObj struct {
	Value Object
}

func (n *ReturnObj) Type() ObjectType { return RETURN_OBJ }

type ListObj struct {
	Values []Object
}

func (n *ListObj) Type() ObjectType { return LIST_OBJ }

type MapObj struct {
	Pairs map[HashKey]KeyValuePair
}

func (n *MapObj) Type() ObjectType { return MAP_OBJ }

type ModuleObj struct {
	Name       string
	ModuleType ModuleType
	Vars       map[string]Object
}

func (n *ModuleObj) Type() ObjectType { return MODULE_OBJ }

type BuiltinObj struct {
	Fn   BuiltinFn
	Name string
}

func (n *BuiltinObj) Type() ObjectType { return BUILTIN_OBJ }

type IteratorObj struct {
	Iterator Iterator
}

func (n *IteratorObj) Type() ObjectType { return ITERATOR_OBJ }

type ErrorObj struct {
	Error error
	Token *token.Token
}

func (n *ErrorObj) Type() ObjectType { return ERROR_OBJ }

// this is gives us a compile time check to see of all the interafaces has been properly implemented
func _() {
	_ = Object(&BooleanObj{})
	_ = Object(&NilObj{})
	_ = Object(&NumberObj{})
	_ = Object(&StringObj{})
	_ = Object(&FnLiteralObj{})
	_ = Object(&ReturnObj{})
	_ = Object(&ListObj{})
	_ = Object(&MapObj{})
	_ = Object(&ModuleObj{})
	_ = Object(&BuiltinObj{})
	_ = Object(&IteratorObj{})
	_ = Object(&ErrorObj{})
}

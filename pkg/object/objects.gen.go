// THIS FILE IS GENERATED. DO NOT EDIT

package object

import "github.com/fredrikkvalvik/temp-lang/pkg/ast"

import "github.com/fredrikkvalvik/temp-lang/pkg/token"

type BooleanObj struct {
	Value bool
}

func (n *BooleanObj) Type() ObjectType { return OBJ_BOOL }

type NilObj struct {
}

func (n *NilObj) Type() ObjectType { return OBJ_NIL }

type NumberObj struct {
	Value float64
}

func (n *NumberObj) Type() ObjectType { return OBJ_NUMBER }

type StringObj struct {
	Value string
}

func (n *StringObj) Type() ObjectType { return OBJ_STRING }

type FnLiteralObj struct {
	Parameters []*ast.IdentifierExpr
	Body       *ast.BlockStmt
	Env        *Environment
}

func (n *FnLiteralObj) Type() ObjectType { return OBJ_FUNCTION_LITERAL }

type ReturnObj struct {
	Value Object
}

func (n *ReturnObj) Type() ObjectType { return OBJ_RETURN }

type ListObj struct {
	Values []Object
}

func (n *ListObj) Type() ObjectType { return OBJ_LIST }

type MapObj struct {
	Pairs map[HashKey]KeyValuePair
}

func (n *MapObj) Type() ObjectType { return OBJ_MAP }

type ModuleObj struct {
	Name       string
	ModuleType ModuleType
	Vars       map[string]Object
}

func (n *ModuleObj) Type() ObjectType { return OBJ_MODULE }

type BuiltinObj struct {
	Fn   BuiltinFn
	Name string
}

func (n *BuiltinObj) Type() ObjectType { return OBJ_BUILTIN }

type IteratorObj struct {
	Iterator Iterator
}

func (n *IteratorObj) Type() ObjectType { return OBJ_ITERATOR }

type ErrorObj struct {
	Error error
	Token *token.Token
}

func (n *ErrorObj) Type() ObjectType { return OBJ_ERROR }

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

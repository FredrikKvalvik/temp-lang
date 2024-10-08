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

type ErrorObj struct {
	Error error
	Token token.Token
}

func (n *ErrorObj) Type() ObjectType { return ERROR_OBJ }

// this is gives us a compile time check to see of all the interafaces has been properly implemented
func typecheck() {
	_ = Object(&BooleanObj{})
	_ = Object(&NilObj{})
	_ = Object(&NumberObj{})
	_ = Object(&StringObj{})
	_ = Object(&FnLiteralObj{})
	_ = Object(&ReturnObj{})
	_ = Object(&ErrorObj{})
}

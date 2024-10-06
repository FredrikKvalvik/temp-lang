package interpreter

import (
	"errors"
	"fmt"

	"github.com/fredrikkvalvik/temp-lang/pkg/ast"
	"github.com/fredrikkvalvik/temp-lang/pkg/object"
	"github.com/fredrikkvalvik/temp-lang/pkg/token"
)

var (
	TypeError            = errors.New("Unexpected type")
	UseOfUndeclaredError = errors.New("Use of undeclared var")
	UnknownNodeError     = errors.New("Unknown node")

	IllegalOperationError     = errors.New("Illegal operation")
	IllegalGlobalReturnError  = errors.New("Illegal return in global scope")
	IllegalRedaclarationError = errors.New("Illegal declaration")
)

// TODO: add line:col numbers to errors

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func illegalOpError(left object.Object, op token.TokenType, right object.Object) *object.Error {
	return &object.Error{
		Error: fmt.Errorf("%w: %s %s %s", IllegalOperationError, left, op, right),
	}
}
func typeMismatchBinaryError(left object.Object, op token.TokenType, right object.Object) *object.Error {
	return &object.Error{Error: fmt.Errorf("%w: %s %s %s", IllegalOperationError, left, op, right)}
}
func typeMismatchUnaryError(op token.TokenType, right object.Object) *object.Error {
	return &object.Error{Error: fmt.Errorf("%w: %s%s", IllegalOperationError, op, right)}
}

// for interal error only. This should only show up in development
func unknownNodeError(node ast.Node) *object.Error {
	return &object.Error{Error: fmt.Errorf("%w: %s", UnknownNodeError, node.Lexeme())}
}

func useOfUnassignVariableError(key string) *object.Error {
	return &object.Error{Error: fmt.Errorf("%w `%s`", UseOfUndeclaredError, key)}
}

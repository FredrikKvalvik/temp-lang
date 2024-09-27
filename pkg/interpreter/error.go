package interpreter

import (
	"fmt"

	"github.com/fredrikkvalvik/temp-lang/pkg/ast"
	"github.com/fredrikkvalvik/temp-lang/pkg/object"
	"github.com/fredrikkvalvik/temp-lang/pkg/token"
)

// TODO: add line:col numbers to errors

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func illegalOpError(left object.Object, op token.TokenType, right object.Object) *object.Error {
	return &object.Error{Message: fmt.Sprintf("Illegal operation: %s %s %s", left, op, right)}
}
func typeMismatchBinaryError(left object.Object, op token.TokenType, right object.Object) *object.Error {
	return &object.Error{Message: fmt.Sprintf("Missmatching type: %s %s %s", left, op, right)}
}
func typeMismatchUnaryError(op token.TokenType, right object.Object) *object.Error {
	return &object.Error{Message: fmt.Sprintf("Missmatching type: %s%s", op, right)}
}

// for interal error only. This should only show up in development
func unknownNodeError(node ast.Node) *object.Error {
	return &object.Error{Message: fmt.Sprintf("Unknown node: %s", node.Lexeme())}
}

func useOfUnassignVariableError(key string) *object.Error {
	return &object.Error{Message: fmt.Sprintf("Use of unassigned var `%s`", key)}
}

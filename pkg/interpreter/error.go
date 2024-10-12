package interpreter

import (
	"errors"
	"fmt"

	"github.com/fredrikkvalvik/temp-lang/pkg/ast"
	"github.com/fredrikkvalvik/temp-lang/pkg/object"
	"github.com/fredrikkvalvik/temp-lang/pkg/token"
)

type InterpreterError error

var (
	TypeError             InterpreterError = errors.New("Unexpected type")
	UseOfUndeclaredError  InterpreterError = errors.New("Use of undeclared var")
	IllegalOperationError InterpreterError = errors.New("Illegal operation")

	IllegalGlobalReturnError  InterpreterError = errors.New("Illegal return in global scope")
	IllegalRedaclarationError InterpreterError = errors.New("Illegal declaration")

	IllegalFloatAsIndexError InterpreterError = errors.New("Can't use decimal as index to list")
	IllegalIndexError        InterpreterError = errors.New("Illegal Index type")
	IndexOutOfBoundsError    InterpreterError = errors.New("Index out of bound")

	// Internal error only
	UnknownNodeError InterpreterError = errors.New("Unknown node")
)

// TODO: add line:col numbers to errors

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func newError(err InterpreterError, msgs ...string) *object.ErrorObj {
	errs := []error{err}
	for _, err := range msgs {
		errs = append(errs, errors.New(err))
	}
	return &object.ErrorObj{Error: errors.Join(errs...)}
}

func illegalOpError(left object.Object, op token.TokenType, right object.Object) *object.ErrorObj {
	return &object.ErrorObj{
		Error: fmt.Errorf("%w: %s %s %s", IllegalOperationError, left, op, right),
	}
}
func typeMismatchBinaryError(left object.Object, op token.TokenType, right object.Object) *object.ErrorObj {
	return &object.ErrorObj{Error: fmt.Errorf("%w: %s %s %s", IllegalOperationError, left, op, right)}
}
func typeMismatchUnaryError(op token.TokenType, right object.Object) *object.ErrorObj {
	return &object.ErrorObj{Error: fmt.Errorf("%w: %s%s", IllegalOperationError, op, right)}
}

// for interal error only. This should only show up in development
func unknownNodeError(node ast.Node) *object.ErrorObj {
	if node == nil {
		return &object.ErrorObj{Error: fmt.Errorf("%w: %s", UnknownNodeError, "nil")}
	} else {
		return &object.ErrorObj{Error: fmt.Errorf("%w: %s", UnknownNodeError, node.String())}
	}
}

func useOfUnassignVariableError(key string) *object.ErrorObj {
	return &object.ErrorObj{Error: fmt.Errorf("%w `%s`", UseOfUndeclaredError, key)}
}

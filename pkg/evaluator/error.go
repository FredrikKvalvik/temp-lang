package evaluator

import (
	"errors"
	"fmt"

	"github.com/fredrikkvalvik/temp-lang/pkg/ast"
	"github.com/fredrikkvalvik/temp-lang/pkg/object"
	"github.com/fredrikkvalvik/temp-lang/pkg/token"
)

type RuntimeError error

var (
	TypeError             RuntimeError = errors.New("Unexpected type")
	UseOfUndeclaredError  RuntimeError = errors.New("Use of undeclared var")
	IllegalOperationError RuntimeError = errors.New("Illegal operation")

	IllegalGlobalReturnError  RuntimeError = errors.New("Illegal return in global scope")
	IllegalRedaclarationError RuntimeError = errors.New("Illegal declaration")
	IllegalAssignmentError    RuntimeError = errors.New("Illegal assignment")

	IllegalFloatAsIndexError RuntimeError = errors.New("Can't use decimal as index to list")
	IllegalIndexError        RuntimeError = errors.New("Illegal Index type")
	IndexOutOfBoundsError    RuntimeError = errors.New("Index out of bound")

	// Internal error only
	UnknownNodeError RuntimeError = errors.New("Unknown node")
)

// TODO: add line:col numbers to errors

func isError(obj object.Object) bool {
	if obj != nil && obj.Type() == object.OBJ_ERROR {
		return true
	}
	return false
}

func newError(err RuntimeError, msgs ...string) *object.ErrorObj {
	errs := []error{err}
	for _, err := range msgs {
		errs = append(errs, errors.New(err))
	}
	return &object.ErrorObj{Error: errors.Join(errs...)}
}

type EnrichErrorParams struct {
	Tok *token.Token
}

// adds more information the the error if possible
func enrichError(err *object.ErrorObj, params *EnrichErrorParams) *object.ErrorObj {
	if err.Token == nil {
		err.Token = params.Tok
	}
	return err
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

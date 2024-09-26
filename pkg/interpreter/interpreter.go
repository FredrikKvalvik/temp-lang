package interpreter

import (
	"fmt"

	"github.com/fredrikkvalvik/temp-lang/pkg/ast"
	"github.com/fredrikkvalvik/temp-lang/pkg/object"
	"github.com/fredrikkvalvik/temp-lang/pkg/token"
)

// SENTINEL VALUES

var TRUE = &object.Boolean{Value: true}
var FALSE = &object.Boolean{Value: false}
var NIL = &object.Nil{}

// main func for interpreter. Recursively evaluate ast and return a value at the end
func Eval(node ast.Node, env *Environment) object.Object {
	// TODO: use assigned value form type
	switch n := node.(type) {
	case *ast.Program:
		return evalProgram(n.Statements, env)
	case *ast.LetStmt:
		key := n.Name.Value
		if env.has(key) {
			return illegalAssignmentError(key)
		}
		value := Eval(n.Value, env)
		env.set(key, value)
		return nil

	case *ast.ExpressionStmt:
		return Eval(n.Expression, env)
	// case *ast.IfStmt:
	// case *ast.BlockStmt:

	// case *ast.UnaryExpr:
	case *ast.BinaryExpr:
		left := Eval(n.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return evalBinaryExpression(left, right, n.Operand)
	case *ast.ParenExpr:
		return Eval(n.Expression, env)

	case *ast.BooleanLiteralExpr:
		return boolObject(n.Value)

	case *ast.NumberLiteralExpr:
		return &object.Number{Value: n.Value}

	case *ast.StringLiteralExpr:
		return &object.String{Value: n.Value}
	default:
		return unknownNodeError(node)
	}
}

func evalBinaryExpression(left, right object.Object, op token.TokenType) object.Object {

	switch {
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringBinaryExpression(left.(*object.String), op, right.(*object.String))

	case left.Type() == object.NUMBER_OBJ && right.Type() == object.NUMBER_OBJ:
		return evalNumberBinaryExpression(left.(*object.Number), op, right.(*object.Number))

	case op == token.EQ:
		// this comparison works because TRUE and FALSE are pointers to singletons
		return boolObject(left == right)
	case op == token.NOT_EQ:
		// this comparison works because TRUE and FALSE are pointers to singletons
		return boolObject(left != right)

	case left.Type() != right.Type():
		return typeMismatchError(left, op, right)
	}

	return illegalOpError(left, op, right)
}
func evalProgram(stmts []ast.Stmt, env *Environment) object.Object {
	var result object.Object
	for _, stmt := range stmts {
		result = Eval(stmt, env)

		if isError(result) {
			return result
		}
	}

	return result
}

// only allow + op on string. all else i illegal
func evalStringBinaryExpression(left *object.String, op token.TokenType, right *object.String) object.Object {
	switch op {
	// string returns
	case token.PLUS:
		return &object.String{Value: left.Value + right.Value}

		// Boolean returns
	case token.EQ:
		return boolObject(left.Value == right.Value)
	case token.NOT_EQ:
		return boolObject(left.Value != right.Value)
	}

	return illegalOpError(left, op, right)
}

// only allow + op on string. all else i illegal
func evalNumberBinaryExpression(left *object.Number, op token.TokenType, right *object.Number) object.Object {

	switch op {
	// Number return
	case token.PLUS:
		return &object.Number{Value: left.Value + right.Value}
	case token.MINUS:
		return &object.Number{Value: left.Value - right.Value}
	case token.SLASH:
		return &object.Number{Value: left.Value / right.Value}
	case token.ASTERISK:
		return &object.Number{Value: left.Value * right.Value}

	// boolean return
	case token.LT:
		return boolObject(left.Value < right.Value)
	case token.GT:
		return boolObject(left.Value > right.Value)
	case token.EQ:
		return boolObject(left.Value == right.Value)
	case token.NOT_EQ:
		return boolObject(left.Value != right.Value)
	}

	return illegalOpError(left, op, right)
}

func illegalOpError(left object.Object, op token.TokenType, right object.Object) *object.Error {
	return &object.Error{Message: fmt.Sprintf("Illegal operation: %s %s %s", left, op, right)}
}
func typeMismatchError(left object.Object, op token.TokenType, right object.Object) *object.Error {
	return &object.Error{Message: fmt.Sprintf("Missmatching type: %s %s %s", left, op, right)}
}
func unknownNodeError(node ast.Node) *object.Error {
	return &object.Error{Message: fmt.Sprintf("Unknown node: %s", node.Lexeme())}
}

func illegalAssignmentError(key string) *object.Error {
	return &object.Error{Message: fmt.Sprintf("Illegal assignment, var %s has already been assign", key)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func boolObject(b bool) *object.Boolean {
	if b {
		return TRUE
	} else {
		return FALSE
	}
}

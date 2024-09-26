package interpreter

import (
	"fmt"

	"github.com/fredrikkvalvik/temp-lang/pkg/ast"
	"github.com/fredrikkvalvik/temp-lang/pkg/object"
	"github.com/fredrikkvalvik/temp-lang/pkg/token"
)

// SENTINEL VALUES
// these should only exist as singleton values. That way,
// we can easily compare the values by pointer

var TRUE = &object.Boolean{Value: true}   // Sentinel value: true
var FALSE = &object.Boolean{Value: false} // Sentinel value: false
var NIL = &object.Nil{}                   // Sentinal value: nil

// main func for interpreter. Recursively evaluate ast and return a value at the end
func Eval(node ast.Node, env *Environment) object.Object {
	// TODO: use assigned value form type
	switch n := node.(type) {
	case *ast.Program:
		return evalProgram(n.Statements, env)
	case *ast.LetStmt:
		key := n.Name.Value
		value := Eval(n.Value, env)
		return env.declareVar(key, value)

	case *ast.ExpressionStmt:
		return Eval(n.Expression, env)

	case *ast.PrintStmt:
		value := Eval(n.Expression, env)
		if !isError(value) {
			fmt.Println(value.Inspect())
		}
		return value
	// case *ast.IfStmt:
	case *ast.BlockStmt:
		return evalBlockStatment(n, env)

	case *ast.UnaryExpr:
		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return evalUnaryExpression(right, n.Operand)
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

	case *ast.IdentifierExpr:
		key := n.Value
		value := env.getVar(key)
		if value == nil {
			return useOfUnassignVariableError(key)
		}
		return value

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

func evalBlockStatment(b *ast.BlockStmt, env *Environment) object.Object {
	scope := NewEnv(env)

	var res object.Object = NIL

	for _, stmt := range b.Statements {
		res = Eval(stmt, scope)
		if isError(res) {
			return res
		}
	}

	return res
}

func evalUnaryExpression(right object.Object, op token.TokenType) object.Object {

	switch {
	case right.Type() == object.NUMBER_OBJ && op == token.MINUS:
		return &object.Number{Value: -right.(*object.Number).Value}

	case right.Type() == object.BOOL_OBJ && op == token.BANG:
		fmt.Print(right)
		if right == TRUE {
			return FALSE
		} else {
			return TRUE
		}
	}

	return typeMismatchUnaryError(op, right)
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
		return typeMismatchBinaryError(left, op, right)
	}

	return illegalOpError(left, op, right)
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

func boolObject(b bool) *object.Boolean {
	if b {
		return TRUE
	} else {
		return FALSE
	}
}

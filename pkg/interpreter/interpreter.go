package interpreter

import (
	"fmt"

	"github.com/fredrikkvalvik/temp-lang/pkg/ast"
	"github.com/fredrikkvalvik/temp-lang/pkg/object"
	"github.com/fredrikkvalvik/temp-lang/pkg/token"
)

var TRUE = &object.Boolean{Value: true}
var FALSE = &object.Boolean{Value: false}
var NIL = &object.Nil{}

type Interpreter struct {
	globalEnv *Environment
	program   *ast.Program

	errors []error
}

func New(program *ast.Program, env *Environment) *Interpreter {
	return &Interpreter{
		program:   program,
		globalEnv: env,
	}
}

func (i *Interpreter) EvalProgram(env *Environment) object.Object {
	var result object.Object
	for _, stmt := range i.program.Statements {
		result = i.Eval(stmt, env)

		if isError(result) {
			return result
		}
	}

	return result
}

// TODO: implement program representation of values
// TODO: implement eval funcs for the different ast.Nodes
func (i *Interpreter) Eval(node ast.Node, env *Environment) object.Object {
	// TODO: use assigned value form type
	switch n := node.(type) {
	// case *ast.LetStmt:
	case *ast.ExpressionStmt:
		return i.Eval(n.Expression, env)
	// case *ast.IfStmt:
	// case *ast.BlockStmt:

	// case *ast.UnaryExpr:
	case *ast.BinaryExpr:
		left := i.Eval(n.Left, env)
		if isError(left) {
			return left
		}
		right := i.Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return i.evalBinaryExpression(left, right, n.Operand)
	case *ast.ParenExpr:
		return i.Eval(n.Expression, env)

	case *ast.BooleanLiteralExpr:
		return boolObject(n.Value)

	case *ast.NumberLiteralExpr:
		return &object.Number{Value: n.Value}

	case *ast.StringLiteralExpr:
		return &object.String{Value: n.Value}

	default:
		i.errors = append(i.errors, fmt.Errorf("unknown node. Could not eval"))
		return unknownNodeError(node)
	}
}

func (i *Interpreter) DidError() bool {
	return len(i.errors) > 0
}
func (i *Interpreter) Errors() []error {
	return i.errors
}

func (i *Interpreter) evalBinaryExpression(left, right object.Object, op token.TokenType) object.Object {

	switch {
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return i.evalStringBinaryExpression(left.(*object.String), op, right.(*object.String))

	case left.Type() == object.NUMBER_OBJ && right.Type() == object.NUMBER_OBJ:
		return i.evalNumberBinaryExpression(left.(*object.Number), op, right.(*object.Number))

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

// only allow + op on string. all else i illegal
func (i *Interpreter) evalStringBinaryExpression(left *object.String, op token.TokenType, right *object.String) object.Object {
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
func (i *Interpreter) evalNumberBinaryExpression(left *object.Number, op token.TokenType, right *object.Number) object.Object {

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

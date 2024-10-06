package interpreter

import (
	"fmt"
	"strings"

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
func Eval(node ast.Node, env *object.Environment) object.Object {
	// TODO: use assigned value form type
	switch n := node.(type) {
	case *ast.Program:
		return evalProgram(n.Statements, env)
	case *ast.LetStmt:
		key := n.Name.Value
		value := Eval(n.Value, env)
		return env.DeclareVar(key, value)

	case *ast.ExpressionStmt:
		return Eval(n.Expression, env)

	case *ast.PrintStmt:
		return evalPrintStatment(n, env)
	case *ast.IfStmt:
		condition := Eval(n.Condition, env)
		if isError(condition) {
			return condition
		}
		if condition == TRUE {
			return Eval(n.Then, env)
		} else if n.Else != nil {
			return Eval(n.Else, env)
		}
		return NIL

	case *ast.BlockStmt:
		return evalBlockStatment(n, env)

	case *ast.EachStmt:
		return evalEachStatment(n, env)

	case *ast.UnaryExpr:
		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return evalUnaryExpression(right, n.Operand)
	case *ast.BinaryExpr:
		// not very clean, but it makes parsing alot easier
		if n.Operand == token.ASSIGN {
			return evalAssignment(n, env)
		}

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
		value := env.GetVar(key)
		if value == nil {
			return useOfUnassignVariableError(key)
		}
		return value

	case *ast.ReturnStmt:
		if env.IsGlobalEnv() {
			return &object.Error{Message: "Illegal return in global scope"}
		}
		// exit early when returning without an expression
		if n.Value == nil {
			return &object.Return{Value: NIL}
		}

		value := Eval(n.Value, env)
		if isError(value) {
			return value
		}
		return &object.Return{Value: value}

	case *ast.FunctionLiteralExpr:
		fn := &object.FnLiteral{
			Parameters: n.Arguments,
			Body:       n.Body,
			Env:        env,
		}
		return fn

	case *ast.CallExpr:
		callee := Eval(n.Callee, env)

		args := evalExpressions(n.Arguments, env)
		if len(args) > 0 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(callee, args)

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

func evalEachStatment(node *ast.EachStmt, env *object.Environment) object.Object {
	scope := object.NewEnv(env)
	if node.Init != nil {
		name := node.Init.Name
		value := Eval(node.Init.Value, scope)
		if isError(value) {
			return value
		}
		scope.DeclareVar(name.Value, value)
	}

	for {
		condition := Eval(node.Condition, scope)
		if condition.Type() != object.BOOL_OBJ {
			return &object.Error{Message: "Condition for loop must evaluate to a boolean value"}
		}

		if condition == FALSE {
			break
		}

		b := Eval(node.Body, scope)
		if isError(b) {
			return b
		}

		if node.Update != nil {
			update := Eval(node.Update, scope)
			if isError(update) {
				return update
			}
		}

	}
	return NIL
}

func evalAssignment(node *ast.BinaryExpr, env *object.Environment) object.Object {
	ident, ok := node.Left.(*ast.IdentifierExpr)
	if !ok {
		return &object.Error{Message: "Can only assign value to identifiers", Token: node.Token}
	}

	val := Eval(node.Right, env)
	if isError(val) {
		return val
	}
	val = env.ReassignVar(ident.Value, val)
	if val == nil {
		return &object.Error{
			Message: fmt.Sprintf("No existing varible with name=%s", ident.Value),
			Token:   node.Token}
	}
	return val
}

func evalPrintStatment(n *ast.PrintStmt, env *object.Environment) object.Object {
	var str strings.Builder
	for idx, expr := range n.Expressions {
		val := Eval(expr, env)
		if isError(val) {
			return val
		}

		str.WriteString(val.Inspect())
		if len(n.Expressions) != idx+1 {
			str.WriteString(", ")
		}
	}
	fmt.Println(str.String())
	return NIL
}

func evalProgram(stmts []ast.Stmt, env *object.Environment) object.Object {
	var result object.Object
	for _, stmt := range stmts {
		result = Eval(stmt, env)

		switch rt := result.(type) {
		case *object.Error:
			return rt

		case *object.Return:
			return rt.Value
		}
	}

	return result
}

func evalBlockStatment(b *ast.BlockStmt, env *object.Environment) object.Object {
	scope := object.NewEnv(env)

	var res object.Object = NIL

	for _, stmt := range b.Statements {
		res = Eval(stmt, scope)
		if isError(res) || res.Type() == object.RETURN_OBJ {
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

func unwrapReturn(obj object.Object) object.Object {
	if ret, ok := obj.(*object.Return); ok {
		return ret.Value
	}
	return obj
}

// On error, will return an error as first, and only item in slice
func evalExpressions(exprs []ast.Expr, env *object.Environment) []object.Object {
	list := make([]object.Object, 0)
	for _, expr := range exprs {
		obj := Eval(expr, env)
		if isError(obj) {
			return []object.Object{obj}
		}
		list = append(list, obj)
	}

	return list
}

func applyFunction(callee object.Object, args []object.Object) object.Object {
	if callee.Type() != object.FUNCTION_LITERAL_OBJ {
		return &object.Error{
			Message: fmt.Sprintf("expected function, got=%s\n", callee.Type()),
		}
	}
	fn := callee.(*object.FnLiteral)
	if len(fn.Parameters) != len(args) {
		return &object.Error{
			Message: fmt.Sprintf("expected number of args=%d, got=%d\n",
				len(fn.Parameters), len(args)),
		}
	}
	scope := object.NewEnv(fn.Env)
	for idx, arg := range fn.Parameters {
		scope.DeclareVar(arg.Value, args[idx])
	}

	evaluated := Eval(fn.Body, scope)
	return unwrapReturn(evaluated)
}

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

var TRUE = &object.BooleanObj{Value: true}   // Sentinel value: true
var FALSE = &object.BooleanObj{Value: false} // Sentinel value: false
var NIL = &object.NilObj{}                   // Sentinal value: nil

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

		// flow
	// -- check for iterable
	// -- if iterable, se if the type is valid
	// -- resolve how to iterate based on the type of iterable
	// -- if name != nil, set the first item to local name var.
	// -- update var at the end of each loop
	// -- each step can be its own function
	case *ast.IterStmt:
		return evalIterStatement(n, env)

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
			return &object.ErrorObj{Error: IllegalGlobalReturnError}
		}
		// exit early when returning without an expression
		if n.Value == nil {
			return &object.ReturnObj{Value: NIL}
		}

		value := Eval(n.Value, env)
		if isError(value) {
			return value
		}
		return &object.ReturnObj{Value: value}

	case *ast.FunctionLiteralExpr:
		fn := &object.FnLiteralObj{
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

	case *ast.ListLiteralExpr:
		list := &object.ListObj{}
		list.Values = evalExpressions(n.Items, env)
		return list

	case *ast.BooleanLiteralExpr:
		return boolObject(n.Value)

	case *ast.NumberLiteralExpr:
		return &object.NumberObj{Value: n.Value}

	case *ast.StringLiteralExpr:
		return &object.StringObj{Value: n.Value}
	default:
		fmt.Printf("%v\n", n)
		return unknownNodeError(node)
	}
}

func evalIterStatement(node *ast.IterStmt, env *object.Environment) object.Object {
	if node.Iterable == nil {
		// when no iterable is found, default to infinite loop
		node.Iterable = &ast.BooleanLiteralExpr{Value: true}
	}

	scope := object.NewEnv(env)
	var name *ast.IdentifierExpr
	if node.Name != nil {
		ident, ok := node.Name.(*ast.IdentifierExpr)
		if !ok {
			return &object.ErrorObj{Error: TypeError, Token: *node.Name.GetToken()}
		}
		scope.DeclareVar(ident.Value, NIL)
		name = ident
	}

	iterable := Eval(node.Iterable, env)
	if isError(iterable) {
		return iterable
	}

	iterator, err := object.NewIterator(iterable)
	if err != nil {
		return err
	}

	var result object.Object = NIL
	for !iterator.Done() {
		val := iterator.Next()
		if name != nil {
			scope.SetVar(name.Value, val)
		}

		result = Eval(node.Body, scope)
		if isError(result) || result.Type() == object.RETURN_OBJ {
			return result
		}
	}

	return result
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

	var ret object.Object = NIL
	for {
		// default condition to true when no condition is defined
		var condition object.Object = TRUE
		if node.Condition != nil {
			Eval(node.Condition, scope)
			if condition.Type() != object.BOOL_OBJ {
				return &object.ErrorObj{Error: fmt.Errorf("Condition for loop must evaluate to a boolean value")}
			}
			if condition == FALSE {
				break
			}
		}

		ret = Eval(node.Body, scope)
		if isError(ret) || ret.Type() == object.RETURN_OBJ {
			return ret
		}

		if node.Update != nil {
			update := Eval(node.Update, scope)
			if isError(update) {
				return update
			}
		}

	}
	return ret
}

func evalAssignment(node *ast.BinaryExpr, env *object.Environment) object.Object {
	ident, ok := node.Left.(*ast.IdentifierExpr)
	if !ok {
		return &object.ErrorObj{Error: fmt.Errorf("Can only assign value to identifiers"), Token: node.Token}
	}

	val := Eval(node.Right, env)
	if isError(val) {
		return val
	}
	val = env.ReassignVar(ident.Value, val)
	if val == nil {
		return &object.ErrorObj{
			Error: UseOfUndeclaredError,
			Token: node.Token}
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
		case *object.ErrorObj:
			return rt

		case *object.ReturnObj:
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
		return &object.NumberObj{Value: -right.(*object.NumberObj).Value}

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
		return evalStringBinaryExpression(left.(*object.StringObj), op, right.(*object.StringObj))

	case left.Type() == object.NUMBER_OBJ && right.Type() == object.NUMBER_OBJ:
		return evalNumberBinaryExpression(left.(*object.NumberObj), op, right.(*object.NumberObj))

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
func evalStringBinaryExpression(left *object.StringObj, op token.TokenType, right *object.StringObj) object.Object {
	switch op {
	// string returns
	case token.PLUS:
		return &object.StringObj{Value: left.Value + right.Value}

		// Boolean returns
	case token.EQ:
		return boolObject(left.Value == right.Value)
	case token.NOT_EQ:
		return boolObject(left.Value != right.Value)
	}

	return illegalOpError(left, op, right)
}

// only allow + op on string. all else i illegal
func evalNumberBinaryExpression(left *object.NumberObj, op token.TokenType, right *object.NumberObj) object.Object {

	switch op {
	// Number return
	case token.PLUS:
		return &object.NumberObj{Value: left.Value + right.Value}
	case token.MINUS:
		return &object.NumberObj{Value: left.Value - right.Value}
	case token.SLASH:
		return &object.NumberObj{Value: left.Value / right.Value}
	case token.ASTERISK:
		return &object.NumberObj{Value: left.Value * right.Value}

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

func boolObject(b bool) *object.BooleanObj {
	if b {
		return TRUE
	} else {
		return FALSE
	}
}

func unwrapReturn(obj object.Object) object.Object {
	if ret, ok := obj.(*object.ReturnObj); ok {
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
		return &object.ErrorObj{
			Error: fmt.Errorf("expected function, got=%s\n", callee.Type()),
		}
	}
	fn := callee.(*object.FnLiteralObj)
	if len(fn.Parameters) != len(args) {
		return &object.ErrorObj{
			Error: fmt.Errorf("expected number of args=%d, got=%d\n",
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

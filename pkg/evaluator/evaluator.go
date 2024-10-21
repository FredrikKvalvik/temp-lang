package evaluator

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

	case *ast.ImportStmt:
		// TODO: implement module logic
		mod := &object.ModuleObj{
			Name:       n.Name.Value,
			ModuleType: object.NATIVE_MODULE,
		}
		return mod

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
		scope := object.NewEnv(env)
		return evalBlockStatment(n, scope)

	case *ast.IterStmt:
		return evalIterStatement(n, env)

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

	case *ast.LogicalExpr:
		left := Eval(n.Left, env)
		if isError(left) {
			return left
		}
		if left.Type() != object.BOOL_OBJ {
			return newError(TypeError, left.Inspect()+" is not of type: boolean")
		}

		if (n.Operand == token.AND && left == TRUE) ||
			(n.Operand == token.OR && left == FALSE) {
			right := Eval(n.Right, env)
			if isError(right) {
				return right
			}

			if right.Type() != object.BOOL_OBJ {
				return newError(TypeError, right.Inspect()+" is not of type: boolean")
			}
			return right
		}

		return left

	case *ast.AssignExpr:
		return evalAssignment(n, env)

	case *ast.IdentifierExpr:
		return evalIdentifier(n, env)

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

	case *ast.IndexExpr:
		left := Eval(n.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(n.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)

	case *ast.ListLiteralExpr:
		list := &object.ListObj{}
		list.Values = evalExpressions(n.Items, env)
		return list

	case *ast.MapLiteralExpr:
		mapLit := &object.MapObj{}
		pairs, err := evalKeyValueExpressions(n.KeyValues, env)
		if err != nil {
			return err
		}
		mapLit.Pairs = pairs

		return mapLit

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

func evalIdentifier(n *ast.IdentifierExpr, env *object.Environment) object.Object {
	// NOTE: we first look for a user-declared variable
	// this is so that if the use declares a variable
	// with the same name as a builtin, we should give
	// the user declared variable priority

	if n.ResolutionDepth > 0 {
		ret := env.GetVar(n.Value, n.ResolutionDepth)
		return ret

	} else if val := env.FindVar(n.Value); val != nil {
		return val

	} else if val, ok := builtins[n.Value]; ok {
		return val
	}

	return newError(UseOfUndeclaredError, n.Value)
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.LIST_OBJ && index.Type() == object.NUMBER_OBJ:
		return evalIndexListExpression(left, index)
	case left.Type() == object.STRING_OBJ && index.Type() == object.NUMBER_OBJ:
		return evalIndexStringExpression(left, index)
	case left.Type() == object.MAP_OBJ:
		return evalIndexMapListExpression(left, index)

	default:
		return newError(TypeError, fmt.Sprintf("%s is not indexeble by %s", left.Inspect(), index.Inspect()))
	}
}

func evalIndexStringExpression(left, index object.Object) object.Object {
	idx := index.(*object.NumberObj).Value
	if !isIntegral(idx) {
		return newError(IllegalFloatAsIndexError)
	}

	str := left.(*object.StringObj).Value
	maxIdx := len(str) - 1

	if int(idx) > maxIdx || idx < 0 {
		return newError(IndexOutOfBoundsError)
	}

	// PERF: extremely inefficient. should look for better solution
	ch := []rune(str)[int(idx)]
	return &object.StringObj{Value: string(ch)}
}

func evalIndexListExpression(left, index object.Object) object.Object {
	idx := index.(*object.NumberObj).Value
	if !isIntegral(idx) {
		return newError(IllegalFloatAsIndexError)
	}

	list := left.(*object.ListObj).Values
	maxIdx := len(list) - 1

	if int(idx) > maxIdx || idx < 0 {
		return newError(IndexOutOfBoundsError)
	}

	return list[int(idx)]
}

func evalIndexMapListExpression(left, index object.Object) object.Object {
	hashable, ok := index.(object.Hashable)
	if !ok {
		return newError(IllegalIndexError, fmt.Sprintf("%s is not a valid key", index.Inspect()))
	}

	hashMap := left.(*object.MapObj).Pairs

	pair, ok := hashMap[hashable.HashKey()]
	if !ok {
		return NIL
	}

	return pair.Value
}

func evalIterStatement(node *ast.IterStmt, env *object.Environment) object.Object {
	if node.Iterable == nil {
		// when no iterable is found, default to infinite loop
		node.Iterable = &ast.BooleanLiteralExpr{Value: true}
	}

	iterable := Eval(node.Iterable, env)
	if isError(iterable) {
		return iterable
	}

	iterator, err := object.NewIterator(iterable)
	if err != nil {
		return err
	}
	var name *ast.IdentifierExpr
	if node.Name != nil {
		ident, ok := node.Name.(*ast.IdentifierExpr)
		if !ok {
			return &object.ErrorObj{Error: TypeError, Token: node.Name.GetToken()}
		}
		name = ident
	}

	var result object.Object = NIL
	for !iterator.Done() {
		val := iterator.Next()

		scope := object.NewEnv(env)

		if name != nil {
			scope.DeclareVar(name.Value, val)
		}

		result = evalBlockStatment(node.Body, scope)
		if isError(result) || result.Type() == object.RETURN_OBJ {
			return result
		}
	}

	return result
}

func evalAssignment(node *ast.AssignExpr, env *object.Environment) object.Object {
	val := Eval(node.Value, env)
	if isError(val) {
		return val
	}

	switch n := node.Assignee.(type) {
	case *ast.IdentifierExpr:
		val = env.ReassignVar(n.Value, val)
		return val

	case *ast.IndexExpr:
		left := Eval(n.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(n.Index, env)
		if isError(index) {
			return index
		}

		return evalIndexAssignment(left, index, val)
	}

	return newError(IllegalAssignmentError)
}

func evalIndexAssignment(assignee, index, value object.Object) object.Object {
	if assignee.Type() == object.LIST_OBJ && index.Type() == object.NUMBER_OBJ {
		return evalIndexListAssignment(index, assignee, value)
	}

	if assignee.Type() == object.MAP_OBJ {
		return evalIndexHashAssignment(index, assignee, value)
	}

	return nil
}

func evalIndexHashAssignment(index object.Object, assignee object.Object, value object.Object) object.Object {
	hash, ok := index.(object.Hashable)
	if !ok {
		return newError(IllegalIndexError)
	}

	pair, ok := assignee.(*object.MapObj).Pairs[hash.HashKey()]
	if ok {
		pair.Value = value
		return pair.Value
	} else {
		kv := object.KeyValuePair{Key: index, Value: value}
		assignee.(*object.MapObj).Pairs[hash.HashKey()] = kv
		return kv.Value
	}

}

func evalIndexListAssignment(index object.Object, assignee object.Object, value object.Object) object.Object {
	idx := index.(*object.NumberObj).Value

	// make sure the index is a whole number
	if !isIntegral(idx) {
		return newError(IllegalFloatAsIndexError)
	}

	// check if index is out of bounds
	if int(idx) >= len(assignee.(*object.ListObj).Values) || int(idx) < 0 {
		return newError(IndexOutOfBoundsError, "check")
	}

	assignee.(*object.ListObj).Values[int(idx)] = value
	return value
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

func evalBlockStatment(b *ast.BlockStmt, scope *object.Environment) object.Object {
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

func evalKeyValueExpressions(expressionMap map[ast.Expr]ast.Expr, env *object.Environment) (
	map[object.HashKey]object.KeyValuePair,
	*object.ErrorObj,
) {
	pairs := make(map[object.HashKey]object.KeyValuePair, len(expressionMap))

	for key, value := range expressionMap {
		k := Eval(key, env)
		if isError(k) {
			return nil, k.(*object.ErrorObj)
		}

		hash, ok := k.(object.Hashable)
		if !ok {
			err := newError(IllegalIndexError, fmt.Sprintf("can't use %s as key", key.String()))
			err.Token = key.GetToken()
			return nil, err
		}

		v := Eval(value, env)
		if isError(k) {
			return nil, v.(*object.ErrorObj)
		}

		pairs[hash.HashKey()] = object.KeyValuePair{Key: k, Value: v}
	}

	return pairs, nil
}

func applyFunction(callee object.Object, args []object.Object) object.Object {

	switch callee.Type() {
	case object.BUILTIN_OBJ:
		val := callee.(*object.BuiltinObj).Fn(args...)
		if val != nil {
			return val
		}
		return NIL

	case object.FUNCTION_LITERAL_OBJ:
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

		evaluated := evalBlockStatment(fn.Body, scope)
		return unwrapReturn(evaluated)

	default:
		return newError(TypeError, fmt.Sprintf("expected function, got=%s\n", callee.Type()))
	}
}

// helper to check if value is a whole number
func isIntegral(val float64) bool {
	return val == float64(int(val))
}

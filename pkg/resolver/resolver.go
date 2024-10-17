// resolver is responsible for doing semantic analysis
// of the program. It checks if returns are valid, clojures
// act as defined, improve lookup time for variables, resolve imports
// by populating the environment
package resolver

import (
	"errors"
	"fmt"
	"slices"

	"github.com/fredrikkvalvik/temp-lang/pkg/ast"
	"github.com/fredrikkvalvik/temp-lang/pkg/object"
)

type ScopeType int

const (
	FunctionScope ScopeType = iota
	IterScope
	// GlobalScope
)

var (
	IllegalDeclarationError           = errors.New("Can't a declare two variables with the same name in the same scope")
	IllegalDefinitionError            = errors.New("Can't define variable that is not declared")
	IllegalUseOfSelfInitError         = errors.New("Can't read local variable in its own initializer")
	IllegalReturnOutsideFunctionError = errors.New("Can't return outside function body")

	UnknownNodeError = errors.New("Unknown node")
)

type Resolver struct {
	scope     Stack[map[string]bool]
	scopeType Stack[ScopeType]
	globalEnv *object.Environment

	Errors []error
}

func New(env *object.Environment) *Resolver {
	r := &Resolver{
		scope:     Stack[map[string]bool]{},
		globalEnv: env,
	}

	return r
}

func (r *Resolver) Resolve(node ast.Node) {
	// fmt.Printf("resolving %T..\n", node)
	switch n := node.(type) {
	case *ast.Program:
		r.enterScope()
		for _, stmt := range n.Statements {
			r.Resolve(stmt)
		}
		r.leaveScope()
		return

	case *ast.BlockStmt:
		r.enterScope()
		for _, stmt := range n.Statements {
			r.Resolve(stmt)
		}
		r.leaveScope()

	case *ast.IfStmt:
		r.Resolve(n.Condition)

		r.Resolve(n.Then)
		if n.Else != nil {
			r.Resolve(n.Else)
		}

	case *ast.IterStmt:
		if n.Name != nil {
			r.Resolve(n.Name)
		}
		if n.Iterable != nil {
			r.Resolve(n.Iterable)
		}
		r.pushScopeType(IterScope)
		r.Resolve(n.Body)
		r.popScopeType()

	case *ast.PrintStmt:
		for _, expr := range n.Expressions {
			r.Resolve(expr)
		}
	case *ast.ExpressionStmt:
		r.Resolve(n.Expression)

	case *ast.ReturnStmt:
		if !r.hasScopeType(FunctionScope) {
			r.Errors = append(r.Errors, IllegalReturnOutsideFunctionError)
		}
		r.Resolve(n.Value)

	case *ast.LetStmt:
		if _, ok := n.Value.(*ast.FunctionLiteralExpr); ok {
			r.declare(n.Name.Value)
			r.define(n.Name.Value)
			r.Resolve(n.Value)
		} else {
			r.declare(n.Name.Value)
			r.Resolve(n.Value)
			r.define(n.Name.Value)
		}

	case *ast.IdentifierExpr:
		if r.scope.IsEmpty() {
			return
		}
		if defined, declared := r.scope.Peek()[n.Value]; declared && !defined {
			r.Errors = append(r.Errors, IllegalUseOfSelfInitError)
		}
		r.resolveLocal(n)

	case *ast.BinaryExpr:
		r.Resolve(n.Left)
		r.Resolve(n.Right)

	case *ast.FunctionLiteralExpr:
		r.enterScope()
		r.pushScopeType(FunctionScope)

		for _, name := range n.Arguments {
			r.declare(name.Value)
			r.define(name.Value)
		}

		r.resolveStmtList(n.Body.Statements)

		r.leaveScope()
		r.popScopeType()

	case *ast.CallExpr:
		r.Resolve(n.Callee)
		r.resolveExprList(n.Arguments)

	case *ast.ListLiteralExpr:
		r.resolveExprList(n.Items)

	case *ast.StringLiteralExpr:
	case *ast.NumberLiteralExpr:
	case *ast.BooleanLiteralExpr:
		// do nothing
	default:
		r.Errors = append(r.Errors, fmt.Errorf("%w: %T", UnknownNodeError, n))

		return
	}
}

func (r *Resolver) resolveLocal(n *ast.IdentifierExpr) {
	for i := r.scope.Size() - 1; i >= 0; i-- {
		if _, ok := r.scope[i][n.Value]; ok {
			n.ResolutionDepth = r.scope.Size() - 1 - i
			return
		}
	}
}

func (r *Resolver) resolveExprList(list []ast.Expr) {
	for _, n := range list {
		r.Resolve(n)
	}
}
func (r *Resolver) resolveStmtList(list []ast.Stmt) {
	for _, n := range list {
		r.Resolve(n)
	}
}

func (r *Resolver) declare(name string) {
	if r.scope.IsEmpty() {
		return
	}
	r.scope.Peek()[name] = false
}
func (r *Resolver) define(name string) {
	if r.scope.IsEmpty() {
		return
	}
	r.scope.Peek()[name] = true
}

// enter a scope with a specific type.
// this will allow ex check the semantics of the program.
//
// example: return is only allowed inside a function body, so to see if that is the case,
// we look through the stack and check if one of the scopeTypes are in fact a function scope
func (r *Resolver) enterScope() {
	r.scope.Push(map[string]bool{})
	// r.scopeType.Push(scopeType)
}

// leave the current scope
func (r *Resolver) leaveScope() {
	r.scope.Pop()
	// r.scopeType.Pop()
}

func (r *Resolver) pushScopeType(st ScopeType) {
	r.scopeType.Push(st)
}
func (r *Resolver) popScopeType() ScopeType {
	return r.scopeType.Pop()
}

func (r *Resolver) hasScopeType(st ScopeType) bool {
	for _, scopeType := range slices.Backward(r.scopeType) {
		if scopeType == st {
			return true
		}
	}

	return false
}
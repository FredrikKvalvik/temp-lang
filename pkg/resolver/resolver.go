// resolver is responsible for doing semantic analysis
// of the program. It checks if returns are valid, clojures
// act as defined, improve lookup time for variables, resolve imports
// by populating the environment
//
// Resolver will also hoist top level function definitions so that we can define
// to allow for calling function before their definition
package resolver

import (
	"errors"
	"fmt"
	"maps"
	"slices"

	"github.com/fredrikkvalvik/temp-lang/pkg/ast"
	"github.com/fredrikkvalvik/temp-lang/pkg/object"
	"github.com/fredrikkvalvik/temp-lang/pkg/std/fmt_std"
	"github.com/fredrikkvalvik/temp-lang/pkg/token"
)

type ScopeType int

const (
	FunctionScope ScopeType = iota
	IterScope
	// GlobalScope
)

var (
	IllegalDeclarationError             = errors.New("Can't a declare two variables with the same name in the same scope")
	IllegalDefinitionError              = errors.New("Can't define variable that is not declared")
	IllegalUseOfSelfInitError           = errors.New("Can't read local variable in its own initializer")
	IllegalReturnOutsideFunctionError   = errors.New("Can't return outside function body")
	IllegalScopedImportError            = errors.New("Can only import in global scope")
	IllegalImportAfterDeclarationsError = errors.New("Can only import at the beginning of the file")

	// error for development. should only be returned when the resolver has not implemented a resolve-case for a node
	UnknownNodeError = errors.New("Resolution for node not implemented")
)

func (r *Resolver) newError(pos token.Pos, err error) {
	r.Errors = append(r.Errors, fmt.Errorf("%s %w", pos, err))
}

// if importPath == stdPath, resolve import to stdLib, else resolve to file path
var stdModules = map[string]*object.ModuleObj{
	"fmt": &fmt_std.Module,
}

type Resolver struct {
	scope     Stack[map[string]bool]
	scopeType Stack[ScopeType]
	globalEnv *object.Environment

	// we are done parsing imports when we resolve any other stmt.
	// imports need to be at the top of the file
	// doneResolvingImports bool

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
	switch n := node.(type) {
	// Program entry point
	case *ast.Program:
		r.resolveImports()
		r.hoistFunctions(n)

		r.resolveStmtList(n.Statements)
		return

	case *ast.ImportStmt:
		// parse all imports.
		if !r.scope.IsEmpty() {
			// r.newError(n.Token.Pos, IllegalScopedImportError)
			r.Errors = append(r.Errors, IllegalScopedImportError)
			return
		}
		// TODO: implement import resolution
		for key := range maps.Keys(stdModules) {
			if key == n.Path {
				r.globalEnv.DeclareVar(key, stdModules[key])
			}
		}
		// if maps.Index(stdModules, n.Path) >= 0 {
		// 	// add logic for

		// }
		// if an import is seen here, that means the the user has imported a module
		// after after other delcarations. that is an error
		// r.newError(n.Token.Pos, IllegalImportAfterDeclarationsError)

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

	case *ast.WhileStmt:
		if n.Condition != nil {
			r.Resolve(n.Condition)
		}
		r.pushScopeType(IterScope)
		r.Resolve(n.Body)
		r.popScopeType()

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
			r.newError(n.Token.Pos, IllegalReturnOutsideFunctionError)
		}
		if n.Value != nil {
			r.Resolve(n.Value)
		}

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
			r.newError(n.Token.Pos, IllegalUseOfSelfInitError)
		}
		r.resolveLocal(n)

	case *ast.BinaryExpr:
		r.Resolve(n.Left)
		r.Resolve(n.Right)

	case *ast.UnaryExpr:
		r.Resolve(n.Right)

	case *ast.ParenExpr:
		r.Resolve(n.Expression)

	case *ast.LogicalExpr:
		r.Resolve(n.Left)
		r.Resolve(n.Right)

	case *ast.AssignExpr:
		r.Resolve(n.Value)
		r.Resolve(n.Assignee)

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

	case *ast.GetExpr:
		r.Resolve(n.Obj)

	case *ast.CallExpr:
		r.Resolve(n.Callee)
		r.resolveExprList(n.Arguments)

	case *ast.ListLiteralExpr:
		r.resolveExprList(n.Items)

	case *ast.MapLiteralExpr:
		for key, value := range n.KeyValues {
			r.Resolve(key)
			r.Resolve(value)
		}

	case *ast.IndexExpr:
		r.Resolve(n.Left)
		r.Resolve(n.Index)

	case *ast.StringLiteralExpr:
	case *ast.NumberLiteralExpr:
	case *ast.BooleanLiteralExpr:
		// do nothing
	default:
		r.Errors = append(r.Errors, fmt.Errorf("%w: %T", UnknownNodeError, n))
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

// filters the function declarations from the program, and moves them to the top of the
// stmt list. this will allow the user call a function being defined later in source
func (r *Resolver) hoistFunctions(program *ast.Program) {
	functions := []ast.Stmt{}
	programStmts := []ast.Stmt{}

	for _, stmt := range program.Statements {

		let, ok := stmt.(*ast.LetStmt)
		if !ok {
			programStmts = append(programStmts, stmt)
			continue
		}

		_, ok = let.Value.(*ast.FunctionLiteralExpr)
		if !ok {
			programStmts = append(programStmts, stmt)
			continue
		}

		// we now know the stmt is a function declaration
		functions = append(functions, let)
	}

	program.Statements = append(functions, programStmts...)
}
func (r *Resolver) resolveImports() {
}

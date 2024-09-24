package interpreter

import (
	"fmt"

	"github.com/fredrikkvalvik/temp-lang/pkg/ast"
	"github.com/fredrikkvalvik/temp-lang/pkg/object"
)

// environment represent the environment of the current execution scope
type Environment struct {
	parent *Environment
	vars   map[string]*object.Object
}

func NewEnv(parent *Environment) *Environment {
	return &Environment{
		parent: parent,
		vars:   make(map[string]*object.Object),
	}
}

type Interpreter struct {
	// TODO: implement representation of environment
	env     *Environment
	program *ast.Program

	errors []error
}

func New(program *ast.Program) *Interpreter {
	return &Interpreter{
		program: program,
	}
}

func (i *Interpreter) EvalProgram() {
	for _, stmt := range i.program.Statements {
		i.Eval(stmt)
	}
}

// TODO: implement program representation of values
// TODO: implement eval funcs for the different ast.Nodes
func (i *Interpreter) Eval(node ast.Node) any {
	// TODO: use assigned value form type
	switch node.(type) {
	// case *ast.LetStmt:
	// case *ast.ExpressionStmt:
	// case *ast.IfStmt:
	// case *ast.BlockStmt:

	// case *ast.UnaryExpr:
	// case *ast.BinaryExpr:
	// case *ast.ParenExpr:
	// case *ast.BooleanLiteralExpr:

	default:
		i.errors = append(i.errors, fmt.Errorf("unknown node. Could not eval"))
		return nil
	}
}

func (i *Interpreter) DidError() bool {
	return len(i.errors) > 0
}
func (i *Interpreter) Errors() []error {
	return i.errors
}

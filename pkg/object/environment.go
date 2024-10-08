package object

import (
	"fmt"
)

// environment represent the environment of the current execution scope
type Environment struct {
	parent *Environment
	vars   map[string]Object
}

func NewEnv(parent *Environment) *Environment {
	return &Environment{
		parent: parent,
		vars:   make(map[string]Object),
	}
}

func (e *Environment) DeclareVar(key string, value Object) Object {
	if e.varInScope(key) {
		return illegalDeclarationError(key)
	}
	e.SetVar(key, value)
	return value
}

func (e *Environment) SetVar(key string, value Object) {
	e.vars[key] = value
}

// walks up the env tree to find the first var with name=key
func (e *Environment) GetVar(key string) Object {
	val, ok := e.vars[key]
	if !ok && e.parent != nil {
		val = e.parent.GetVar(key)
	}

	return val
}

// move up the parent tree and look for a variable to reassign. return nil if none are found.
func (e *Environment) ReassignVar(key string, value Object) Object {
	// three cases can happen
	// - variable is found i current scope. we set the value and return it
	// - value is not found i current scope, but scope is not global. we move to parent scope and try again
	// - value is not found i current scope and scope is global. this means we didnt find a variable, so we return nil
	if e.varInScope(key) {
		e.SetVar(key, value)
		return value

	} else if !e.varInScope(key) && !e.IsGlobalEnv() {
		return e.parent.ReassignVar(key, value)

	} else {
		return nil
	}
}

// checks the current scope for an existing variable name
func (e *Environment) varInScope(key string) bool {
	_, ok := e.vars[key]
	return ok
}

// helper for checking if the current  scope is global
func (e *Environment) IsGlobalEnv() bool {
	return e.parent == nil
}

func illegalDeclarationError(key string) *ErrorObj {
	return &ErrorObj{Error: fmt.Errorf(
		"Illegal declaration, var `%s` has already been declared in this scope",
		key,
	)}
}

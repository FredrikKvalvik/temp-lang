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
	if e.hasVar(key) {
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

func (e *Environment) hasVar(key string) bool {
	_, ok := e.vars[key]
	return ok
}

// might be useful?
func (e *Environment) isGlobalEnv() bool {
	return e.parent == nil
}

func illegalDeclarationError(key string) *Error {
	return &Error{Message: fmt.Sprintf("Illegal declaration, var `%s` has already been declared in this scope", key)}
}

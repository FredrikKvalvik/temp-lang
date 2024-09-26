package interpreter

import "github.com/fredrikkvalvik/temp-lang/pkg/object"

// environment represent the environment of the current execution scope
type Environment struct {
	parent *Environment
	vars   map[string]object.Object
}

func NewEnv(parent *Environment) *Environment {
	return &Environment{
		parent: parent,
		vars:   make(map[string]object.Object),
	}
}

func (e *Environment) declareVar(key string, value object.Object) object.Object {
	if e.hasVar(key) {
		return illegalDeclarationError(key)
	}
	e.setVar(key, value)
	return NIL
}

func (e *Environment) setVar(key string, value object.Object) {
	e.vars[key] = value
}

// walks up the env tree to find the first var with name=key
func (e *Environment) getVar(key string) object.Object {
	val, ok := e.vars[key]
	if !ok && e.parent != nil {
		val = e.parent.getVar(key)
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

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

func (e *Environment) set(key string, value object.Object) {
	e.vars[key] = value
}

func (e *Environment) get(key string) object.Object {
	val, ok := e.vars[key]
	if !ok && e.parent != nil {
		val = e.parent.get(key)
	}

	return val
}

func (e *Environment) has(key string) bool {
	_, ok := e.vars[key]
	return ok
}

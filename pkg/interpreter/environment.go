package interpreter

import "github.com/fredrikkvalvik/temp-lang/pkg/object"

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
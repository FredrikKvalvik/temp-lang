package iterstd

import (
	"github.com/fredrikkvalvik/temp-lang/pkg/object"
)

var Module = object.ModuleObj{
	Name:       "iter",
	ModuleType: object.NATIVE_MODULE,
	Vars:       vars,
}
var vars = map[string]object.Object{
	// "range": &object.BuiltinObj{
	// 	Name: "range",
	// 	Fn: func(args ...object.Object) object.Object {

	// 	},
	// },
}

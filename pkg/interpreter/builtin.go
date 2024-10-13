package interpreter

import "github.com/fredrikkvalvik/temp-lang/pkg/object"

var builtins = map[string]*object.BuiltinObj{
	"len": &object.BuiltinObj{Name: "len", Fn: object.LenBuiltin},
}

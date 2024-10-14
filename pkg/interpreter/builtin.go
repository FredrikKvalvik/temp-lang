package interpreter

import "github.com/fredrikkvalvik/temp-lang/pkg/object"

var builtins = map[string]*object.BuiltinObj{
	"len":  {Name: "len", Fn: object.LenBuiltin},
	"push": {Name: "push", Fn: object.PushBuiltin},
	"pop":  {Name: "pop", Fn: object.PopBuiltin},
	"str":  {Name: "str", Fn: object.StrBuiltin},
}

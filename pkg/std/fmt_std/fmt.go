package fmt_std

import (
	"fmt"
	"os"
	"strings"

	"github.com/fredrikkvalvik/temp-lang/pkg/object"
)

var Module = object.ModuleObj{
	Name:       "fmt",
	ModuleType: object.NATIVE_MODULE,
	Vars:       vars,
}
var vars = map[string]object.Object{
	"println": &object.BuiltinObj{
		Name: "println",
		Fn: func(args ...object.Object) object.Object {
			var str strings.Builder

			str.WriteString(objectsToString(args...))
			str.WriteString("\n")

			fmt.Fprint(os.Stdout, str.String())
			return nil
		},
	},

	"string": &object.BuiltinObj{
		Name: "string",
		Fn: func(args ...object.Object) object.Object {
			strObj := &object.StringObj{}

			var str strings.Builder
			str.WriteString(objectsToString(args...))

			strObj.Value = str.String()

			return strObj
		},
	},
}

func toString(obj object.Object) string {

	switch v := obj.(type) {
	case *object.NumberObj:
		return fmt.Sprint(v.Value)
	case *object.StringObj:
		return fmt.Sprint(v.Value)
	case *object.BooleanObj:
		return fmt.Sprintf("%v", v.Value)
	case *object.NilObj:
		return "nil"

	default:
		return v.Inspect()
	}
}

func objectsToString(objs ...object.Object) string {

	var str strings.Builder

	for idx, arg := range objs {
		str.WriteString(toString(arg))
		if idx != len(objs)-1 {
			str.WriteString(" ")
		}
	}

	return str.String()
}

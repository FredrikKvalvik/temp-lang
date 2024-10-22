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

			for _, arg := range args {
				str.WriteString(toString(arg))
				str.WriteString(" ")
			}
			str.WriteString("\n")

			fmt.Fprint(os.Stdout, str.String())
			return nil
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

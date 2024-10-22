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
	"print": &object.BuiltinObj{
		Name: "print",
		Fn: func(args ...object.Object) object.Object {
			var str strings.Builder

			for _, arg := range args {
				switch v := arg.(type) {
				case *object.NumberObj:
					fmt.Fprint(&str, v.Value)
					// TODO: implement print for each type of value

				default:
					fmt.Fprint(&str, v.Inspect())
				}
			}
			fmt.Fprint(os.Stdout, str)
			return nil
		},
	},
}

package object

import (
	"errors"
	"fmt"
)

var (
	ArityError = errors.New("wrong number of args")
	TypeError  = errors.New("invalid type")
)

// takes 1 and returns the length of the object.
// will return nil if there is no way to return a length
func LenBuiltin(args ...Object) Object {
	if len(args) != 1 {
		return &ErrorObj{Error: fmt.Errorf("%w: expected=%d, got=%d", ArityError, 1, len(args))}
	}

	arg := args[0]

	switch argument := arg.(type) {
	case *StringObj:
		return &NumberObj{Value: float64(len([]rune(argument.Value)))}
	case *ListObj:
		return &NumberObj{Value: float64(len(argument.Values))}
	case *MapObj:
		return &NumberObj{Value: float64(len(argument.Pairs))}
	}

	return nil
}

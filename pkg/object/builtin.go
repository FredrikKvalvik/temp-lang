package object

import (
	"errors"
	"fmt"
)

var (
	ArityError = errors.New("wrong number of args")
	TypeError  = errors.New("invalid type")
)

type ErrorBuf struct {
	Err *ErrorObj
}

func (e *ErrorBuf) Run(errFn func() *ErrorObj) {
	if e.Err == nil {
		e.Err = errFn()
	}
}

// util for builtin functions to return an error if check evaluates to false. return nil on ok
func CheckArity(args []Object, arity int) *ErrorObj {
	if len(args) != arity {
		return &ErrorObj{Error: fmt.Errorf("%w: expected %d args, got %d", ArityError, arity, len(args))}
	}
	return nil
}

// util for builtin functions to return an error if check evaluates to false. return nil on ok
func CheckObjectType(obj Object, typ ObjectType) *ErrorObj {
	if obj.Type() != typ {
		return &ErrorObj{Error: fmt.Errorf("%w: expected %s, got %s", TypeError, typ, obj.Type())}
	}
	return nil
}

// util for builtin functions to return an error if check evaluates to false. return nil on ok
func CheckIntegral(n float64) *ErrorObj {
	if !isIntegral(n) {
		return &ErrorObj{Error: fmt.Errorf("%w: expected integer, but %v is float", TypeError, n)}
	}
	return nil
}

// Arity: 1
//
// Arg0: list | map | string
//
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

// Arity: >1
//
// Arg0: list, Arg>0: any
//
// push appends items to the end of a list
// first arg must be type=list, every following argument will be pushed
// to the list in the order they are gotten.
// return the reference to the list
func PushBuiltin(args ...Object) Object {
	if len(args) < 2 {
		return &ErrorObj{Error: fmt.Errorf("%w: expected target list and item(s), got=%d", ArityError, len(args))}
	}
	list := args[0]
	if list.Type() != OBJ_LIST {
		return &ErrorObj{Error: fmt.Errorf("%w: expected list, got=%s", TypeError, list.Type())}
	}

	items := args[1:]
	list.(*ListObj).Values = append(list.(*ListObj).Values, items...)

	return list
}

// Arity: 1
//
// Arg0: list
//
// pop removes the last element from a list and returns it.
// if pop is used on an empty list, return nil
func PopBuiltin(args ...Object) Object {
	var ebuf ErrorBuf

	ebuf.Run(func() *ErrorObj { return CheckArity(args, 1) })
	ebuf.Run(func() *ErrorObj { return CheckObjectType(args[0], OBJ_LIST) })
	if ebuf.Err != nil {
		return ebuf.Err
	}

	arg := args[0]

	list := arg.(*ListObj)

	if len(list.Values) == 0 {
		return nil
	}

	last := list.Values[len(list.Values)-1]
	list.Values = list.Values[:len(list.Values)-1]

	return last
}

// Arity: 1
//
// Arg0: any
//
// pop removes the last element from a list and returns it.
// if pop is used on an empty list, return nil
func StrBuiltin(args ...Object) Object {
	if err := CheckArity(args, 1); err != nil {
		return err
	}
	arg := args[0]

	return &StringObj{arg.Inspect()}
}

// Arity: 3
//
//	Arg0: number
//	Arg1: number
//	Arg2: number
//
// pop removes the last element from a list and returns it.
// if pop is used on an empty list, return nil
func RangeBuiltin(args ...Object) Object {
	var ebuf ErrorBuf

	ebuf.Run(func() *ErrorObj { return CheckArity(args, 3) })
	ebuf.Run(func() *ErrorObj { return CheckObjectType(args[0], OBJ_NUMBER) })
	ebuf.Run(func() *ErrorObj { return CheckObjectType(args[1], OBJ_NUMBER) })
	ebuf.Run(func() *ErrorObj { return CheckObjectType(args[2], OBJ_NUMBER) })

	if ebuf.Err != nil {
		return ebuf.Err
	}

	startObj := args[0].(*NumberObj)
	endObj := args[1].(*NumberObj)
	stepObj := args[2].(*NumberObj)

	ebuf.Run(func() *ErrorObj { return CheckIntegral(startObj.Value) })
	ebuf.Run(func() *ErrorObj { return CheckIntegral(endObj.Value) })
	ebuf.Run(func() *ErrorObj { return CheckIntegral(stepObj.Value) })
	ebuf.Run(func() *ErrorObj {
		if stepObj.Value >= 0 {
			return &ErrorObj{Error: fmt.Errorf("%w: step value must be a none-zero, positive number", TypeError)}
		}
		return nil
	})

	if ebuf.Err != nil {
		return ebuf.Err
	}

	start := int(startObj.Value)
	end := int(endObj.Value)
	step := int(stepObj.Value)
	isNegative := start > end

	if isNegative {
		step = -step
	}

	return &IteratorObj{
		Iterator: newRangeIterator(start, end, step),
	}
}

// Arity: 1
//
// Arg0: any
//
// returns an iterator based on the argument
func IterBuiltin(args ...Object) Object {
	if len(args) != 1 {
		return &ErrorObj{Error: fmt.Errorf("%w: expected=%d, got=%d", ArityError, 1, len(args))}
	}
	arg := args[0]

	iterator, err := NewIterator(arg)
	if err != nil {
		return err
	}

	return &IteratorObj{
		Iterator: iterator,
	}
}

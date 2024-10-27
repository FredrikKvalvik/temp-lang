package object

import (
	"errors"
	"fmt"
)

var (
	ArityError = errors.New("wrong number of args")
	TypeError  = errors.New("invalid type")
)

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
	if len(args) != 1 {
		return &ErrorObj{Error: fmt.Errorf("%w: expected=%d, got=%d", ArityError, 1, len(args))}
	}
	arg := args[0]
	if arg.Type() != OBJ_LIST {
		return &ErrorObj{Error: fmt.Errorf("%w: expected list, got=%s", TypeError, arg.Type())}
	}

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
	if len(args) != 1 {
		return &ErrorObj{Error: fmt.Errorf("%w: expected=%d, got=%d", ArityError, 1, len(args))}
	}
	arg := args[0]

	return &StringObj{arg.Inspect()}
}

// Arity: 2-3
//
//	Arg0: number
//	Arg1: number
//	Arg2: number
//
// pop removes the last element from a list and returns it.
// if pop is used on an empty list, return nil
func RangeBuiltin(args ...Object) Object {
	if len(args) < 2 || len(args) > 3 {
		return &ErrorObj{Error: fmt.Errorf("%w: expected=2-3, got=%d", ArityError, len(args))}
	}

	startObj, ok := args[0].(*NumberObj)
	if !ok {
		return &ErrorObj{Error: fmt.Errorf("%w: expected number, got=%s", TypeError, args[0].Type())}
	}
	endObj, ok := args[1].(*NumberObj)
	if !ok {
		return &ErrorObj{Error: fmt.Errorf("%w: expected number, got=%s", TypeError, args[1].Type())}
	}

	if !isIntegral(startObj.Value) {
		return &ErrorObj{Error: fmt.Errorf("%w: expected number to be integer, got=%v", TypeError, startObj.Value)}
	}
	if !isIntegral(endObj.Value) {
		return &ErrorObj{Error: fmt.Errorf("%w: expected number to be integer, got=%v", TypeError, endObj.Value)}
	}

	start := int(startObj.Value)
	end := int(endObj.Value)
	step := 1
	isNegative := start > end

	if len(args) == 3 {
		stepObj, ok := args[2].(*NumberObj)
		if !ok {
			return &ErrorObj{Error: fmt.Errorf("%w: expected number, got=%s", TypeError, args[2].Type())}
		}
		if !isIntegral(stepObj.Value) {
			return &ErrorObj{Error: fmt.Errorf("%w: expected number to be integer, got=%v", TypeError, stepObj.Value)}
		}
		if step < 0 {
			return &ErrorObj{Error: fmt.Errorf("%w: step value must be a none-zero, positive number", TypeError)}
		}

		step = int(stepObj.Value)
	}
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

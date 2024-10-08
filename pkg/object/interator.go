package object

import (
	"fmt"
	"strings"
)

type IteratorType int

const (
	NUMBER_ITER IteratorType = iota
	STRING_ITER
)

// return true while iterator is returning items. return false when end of loop is finisjed
type Iterator interface {
	Next() Object
	Done() bool
}

func NewIterator(iterable Object) (Iterator, *ErrorObj) {
	switch it := iterable.(type) {
	case *StringObj:
		return newStringIterator(it), nil
	case *NumberObj:
		return newNumberIterator(it), nil
	case *BooleanObj:
		return newBooleanIterator(it), nil

	default:
		return nil, &ErrorObj{Error: fmt.Errorf("%s is not iterable", iterable.Type())}
	}

}

type StringIter struct {
	reader *strings.Reader
}

func newStringIterator(str *StringObj) *StringIter {
	r := strings.NewReader(str.Value)
	return &StringIter{
		reader: r,
	}
}

func (si *StringIter) Next() Object {
	if si.Done() {
		return nil
	}

	ch, _, err := si.reader.ReadRune()

	if err != nil {
		return &ErrorObj{Error: fmt.Errorf("Could not read string")}
	}
	str := &StringObj{Value: string(ch)}

	return str
}

func (si *StringIter) Done() bool { return si.reader.Len() == 0 }

type NumberIter struct {
	number *NumberObj
	index  int
}

func newNumberIterator(num *NumberObj) *NumberIter {
	return &NumberIter{
		number: num,
	}
}
func (ni *NumberIter) Next() Object {
	n := &NumberObj{Value: float64(ni.index)}
	ni.index += 1
	return n
}
func (ni *NumberIter) Done() bool { return ni.index >= int(ni.number.Value) }

type BooleanIter struct {
	bool *BooleanObj
}

func newBooleanIterator(bool *BooleanObj) *BooleanIter {
	return &BooleanIter{
		bool: bool,
	}
}
func (ni *BooleanIter) Next() Object {
	return ni.bool
}
func (ni *BooleanIter) Done() bool { return !ni.bool.Value }

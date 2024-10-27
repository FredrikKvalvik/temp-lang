//go:generate go run golang.org/x/tools/cmd/stringer -type=IteratorType
package object

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

type IteratorType int

const (
	ITER_NUMBER IteratorType = iota
	ITER_STRING
	ITER_LIST
	ITER_MAP
	ITER_RANGE
)

// return true while iterator is returning items. return false when end of loop is finisjed
type Iterator interface {
	Next() Object
	Done() bool
	Type() IteratorType
}

func NewIterator(iterable Object) (Iterator, *ErrorObj) {
	switch it := iterable.(type) {
	case *StringObj:
		return newStringIterator(it), nil
	case *NumberObj:
		if !isIntegral(it.Value) {
			return nil, &ErrorObj{Error: fmt.Errorf("Number iterator must be whole number, got=%v", it.Value)}
		}
		return newNumberIterator(it), nil
	case *ListObj:
		return newListIterator(it), nil
	case *MapObj:
		return newMapIterator(it), nil

	// used for arbitrary iterators
	case *IteratorObj:
		return it.Iterator, nil

	default:
		return nil, &ErrorObj{Error: fmt.Errorf("%s is not iterable", iterable.Type())}
	}

}

// STRING_ITER
type StringIter struct {
	reader *strings.Reader
}

func newStringIterator(str *StringObj) *StringIter {
	r := strings.NewReader(str.Value)
	return &StringIter{reader: r}
}

func (si *StringIter) Type() IteratorType { return ITER_STRING }

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

// NUMBER_ITER

type NumberIter struct {
	number *NumberObj
	index  int
}

func newNumberIterator(num *NumberObj) *NumberIter {

	return &NumberIter{
		number: num,
	}
}
func (i *NumberIter) Type() IteratorType { return ITER_NUMBER }
func (ni *NumberIter) Next() Object {

	n := &NumberObj{Value: float64(ni.index)}
	ni.index += 1
	return n
}
func (ni *NumberIter) Done() bool { return ni.index >= int(ni.number.Value) }

// LIST_ITER

type ListIter struct {
	values []Object
	idx    int
}

func newListIterator(list *ListObj) *ListIter {
	return &ListIter{
		values: list.Values,
	}
}
func (i *ListIter) Type() IteratorType { return ITER_LIST }
func (li *ListIter) Next() Object {
	value := li.values[li.idx]
	li.idx++
	return value
}
func (li *ListIter) Done() bool { return li.idx >= len(li.values) }

// helper to check if value is a whole number

// MAP_ITER

// iterate over the keys of a map
type MapIter struct {
	pairs []KeyValuePair
	idx   int
}

func newMapIterator(m *MapObj) *MapIter {
	pairs := slices.Collect(maps.Values(m.Pairs))
	return &MapIter{
		pairs: pairs,
	}
}

func (mi *MapIter) Type() IteratorType { return ITER_MAP }

func (mi *MapIter) Next() Object {
	value := mi.pairs[mi.idx]
	mi.idx++
	return value.Key
}

func (mi *MapIter) Done() bool { return mi.idx >= len(mi.pairs) }

// helper to check if value is a whole number
func isIntegral(val float64) bool {
	return val == float64(int(val))
}

// Range

type RangeIter struct {
	start int // iterator starts
	end   int // iterator ends
	step  int // increment for the loop

	index int // current index
}

func newRangeIterator(start, end, step int) *RangeIter {
	return &RangeIter{
		start: start,
		end:   end,
		step:  step,

		index: start,
	}
}
func (ri *RangeIter) Type() IteratorType { return ITER_RANGE }
func (ri *RangeIter) Next() Object {
	n := &NumberObj{Value: float64(ri.index)}
	ri.index += ri.step
	return n
}
func (ri *RangeIter) Done() bool {
	if ri.start > ri.end {
		// means we are iterating in negative direction
		return ri.index <= int(ri.end)

	} else {
		return ri.index >= int(ri.end)
	}
}

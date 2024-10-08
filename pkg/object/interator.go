package object

type Iterator interface {
	Next() Object
	Done() bool
}

func NewIterator(iterable Object) Iterator {
	switch it := iterable.(type) {
	case *String:
		return &StringIterator{
			string: it,
		}
	case *Number:
		return &NumberIterator{
			number: it,
		}
	}

	panic("no iterator implemented for type")
}

type StringIterator struct {
	string *String
	index  int
}

func (si *StringIterator) Next() Object { return nil }
func (si *StringIterator) Done() bool   { return true }

type NumberIterator struct {
	number *Number
	index  int
}

func (ni *NumberIterator) Next() Object { return nil }
func (ni *NumberIterator) Done() bool   { return true }

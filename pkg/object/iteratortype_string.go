// Code generated by "stringer -type=IteratorType"; DO NOT EDIT.

package object

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ITER_NUMBER-0]
	_ = x[ITER_STRING-1]
	_ = x[ITER_LIST-2]
	_ = x[ITER_MAP-3]
	_ = x[ITER_RANGE-4]
}

const _IteratorType_name = "ITER_NUMBERITER_STRINGITER_LISTITER_MAPITER_RANGE"

var _IteratorType_index = [...]uint8{0, 11, 22, 31, 39, 49}

func (i IteratorType) String() string {
	if i < 0 || i >= IteratorType(len(_IteratorType_index)-1) {
		return "IteratorType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _IteratorType_name[_IteratorType_index[i]:_IteratorType_index[i+1]]
}

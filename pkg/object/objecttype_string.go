// Code generated by "stringer -type=ObjectType"; DO NOT EDIT.

package object

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[BOOL_OBJ-0]
	_ = x[NULL_OBJ-1]
	_ = x[NUMBER_OBJ-2]
	_ = x[STRING_OBJ-3]
}

const _ObjectType_name = "BOOL_OBJNULL_OBJNUMBER_OBJSTRING_OBJ"

var _ObjectType_index = [...]uint8{0, 8, 16, 26, 36}

func (i ObjectType) String() string {
	if i < 0 || i >= ObjectType(len(_ObjectType_index)-1) {
		return "ObjectType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ObjectType_name[_ObjectType_index[i]:_ObjectType_index[i+1]]
}

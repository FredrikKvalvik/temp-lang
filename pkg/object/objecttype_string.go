// Code generated by "stringer -type=ObjectType"; DO NOT EDIT.

package object

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[BOOL_OBJ-1]
	_ = x[NIL_OBJ-2]
	_ = x[NUMBER_OBJ-3]
	_ = x[STRING_OBJ-4]
	_ = x[FUNCTION_LITERAL_OBJ-5]
	_ = x[RETURN_OBJ-6]
	_ = x[LIST_OBJ-7]
	_ = x[MAP_OBJ-8]
	_ = x[BUILTIN_OBJ-9]
	_ = x[ERROR_OBJ-10]
}

const _ObjectType_name = "BOOL_OBJNIL_OBJNUMBER_OBJSTRING_OBJFUNCTION_LITERAL_OBJRETURN_OBJLIST_OBJMAP_OBJBUILTIN_OBJERROR_OBJ"

var _ObjectType_index = [...]uint8{0, 8, 15, 25, 35, 55, 65, 73, 80, 91, 100}

func (i ObjectType) String() string {
	i -= 1
	if i < 0 || i >= ObjectType(len(_ObjectType_index)-1) {
		return "ObjectType(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _ObjectType_name[_ObjectType_index[i]:_ObjectType_index[i+1]]
}

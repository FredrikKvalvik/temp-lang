// Code generated by "stringer -type=TokenType"; DO NOT EDIT.

package token

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ILLEGAL-0]
	_ = x[EOF-1]
	_ = x[IDENT-2]
	_ = x[NUMBER-3]
	_ = x[STRING-4]
	_ = x[ASSIGN-5]
	_ = x[PLUS-6]
	_ = x[MINUS-7]
	_ = x[BANG-8]
	_ = x[ASTERISK-9]
	_ = x[SLASH-10]
	_ = x[EQ-11]
	_ = x[NOT_EQ-12]
	_ = x[LT-13]
	_ = x[GT-14]
	_ = x[AND-15]
	_ = x[OR-16]
	_ = x[COMMA-17]
	_ = x[DOT-18]
	_ = x[SEMICOLON-19]
	_ = x[COLON-20]
	_ = x[LPAREN-21]
	_ = x[RPAREN-22]
	_ = x[LBRACE-23]
	_ = x[RBRACE-24]
	_ = x[LBRACKET-25]
	_ = x[RBRACKET-26]
	_ = x[FUNCTION-27]
	_ = x[LET-28]
	_ = x[TRUE-29]
	_ = x[FALSE-30]
	_ = x[IF-31]
	_ = x[ELSE-32]
	_ = x[RETURN-33]
	_ = x[PRINT-34]
}

const _TokenType_name = "ILLEGALEOFIDENTNUMBERSTRINGASSIGNPLUSMINUSBANGASTERISKSLASHEQNOT_EQLTGTANDORCOMMADOTSEMICOLONCOLONLPARENRPARENLBRACERBRACELBRACKETRBRACKETFUNCTIONLETTRUEFALSEIFELSERETURNPRINT"

var _TokenType_index = [...]uint8{0, 7, 10, 15, 21, 27, 33, 37, 42, 46, 54, 59, 61, 67, 69, 71, 74, 76, 81, 84, 93, 98, 104, 110, 116, 122, 130, 138, 146, 149, 153, 158, 160, 164, 170, 175}

func (i TokenType) String() string {
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}

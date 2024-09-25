package interpreter

import (
	"fmt"
	"testing"

	"github.com/fredrikkvalvik/temp-lang/pkg/lexer"
	"github.com/fredrikkvalvik/temp-lang/pkg/object"
	"github.com/fredrikkvalvik/temp-lang/pkg/parser"
	"github.com/fredrikkvalvik/temp-lang/pkg/tester"
)

func TestBinaryExpression(t *testing.T) {
	tr := tester.New(t, "")

	tests := []struct {
		input        string
		expectedVal  any
		expectedType object.ObjectType
	}{
		// number returns
		{"2+2",
			float64(4), object.NUMBER_OBJ},
		{"2-2",
			float64(0), object.NUMBER_OBJ},
		{"10 / 2",
			float64(5), object.NUMBER_OBJ},
		{"10 * 2",
			float64(20), object.NUMBER_OBJ},
		{"10 + 2 * 100",
			float64(210), object.NUMBER_OBJ},

		// boolean returns
		{"10 == 2",
			false, object.BOOL_OBJ},
		{"10 != 2",
			true, object.BOOL_OBJ},
		{`"hello" == "hello"`,
			true, object.BOOL_OBJ},
		{`"hello" != "goodbye"`,
			true, object.BOOL_OBJ},
		{"10 < 2",
			false, object.BOOL_OBJ},
		{"10 > 2",
			true, object.BOOL_OBJ},

		// string returns
		{`"hello" + " " + "world"`,
			"hello world", object.STRING_OBJ},

		// error returns
		{`"hello" - " world"`,
			nil, object.ERROR_OBJ},
	}

	for index, tt := range tests {
		result, _ := testEvalProgram(t, tt.input)

		tr.SetName(fmt.Sprint(index))

		tr.AssertEqual(result.Type(), tt.expectedType)

		switch result.Type() {
		case object.NUMBER_OBJ:
			tr.AssertEqual(result.(*object.Number).Value, tt.expectedVal)
		case object.STRING_OBJ:
			tr.AssertEqual(result.(*object.String).Value, tt.expectedVal)
		case object.BOOL_OBJ:
			tr.AssertEqual(result.(*object.Boolean).Value, tt.expectedVal)
		case object.ERROR_OBJ:
			// TODO: check error message
		default:
			tr.T.Errorf("Unexpected object type=%s\n", result.Type())
		}
	}
}

func testEvalProgram(t *testing.T, input string) (object.Object, *Environment) {
	t.Helper()
	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	if p.DidError() {
		for _, err := range p.Errors() {
			t.Error(err, "\n")
		}
		return nil, nil
	}

	env := NewEnv(nil)

	return Eval(program, env), env
}

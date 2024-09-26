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
		tr.SetName(fmt.Sprint(index))

		result, _ := testEvalProgram(tr, tt.input)

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

func TestLetStatement(t *testing.T) {

	input := "let ident = 10"

	tr := tester.New(t, input)

	res, e := testEvalProgram(tr, input)

	value := e.get("ident")

	tr.AssertNil(res)
	tr.AssertNotNil(value)
	tr.AssertEqual(value.Type(), object.NUMBER_OBJ)
	tr.AssertEqual(value.(*object.Number).Value, float64(10))
}

func testEvalProgram(tr *tester.Tester, input string) (object.Object, *Environment) {
	tr.T.Helper()

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	if p.DidError() {
		for _, err := range p.Errors() {
			tr.T.Error(err, "\n")
		}
		return nil, nil
	}

	env := NewEnv(nil)

	return Eval(program, env), env
}

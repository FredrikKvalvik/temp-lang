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

func TestUnaryExpressions(t *testing.T) {
	tr := tester.New(t, "")

	tests := []struct {
		input         string
		expectedType  object.ObjectType
		expectedValue any
	}{
		{"!false",
			object.BOOL_OBJ, true},
		{"!true",
			object.BOOL_OBJ, false},
		{"!!true",
			object.BOOL_OBJ, true},
		{"-10",
			object.NUMBER_OBJ, float64(-10)},
		{"--10",
			object.NUMBER_OBJ, float64(10)},
		{"-true",
			object.ERROR_OBJ, nil},
	}

	for idx, tt := range tests {
		tr.SetName(fmt.Sprintf("[%d]", idx))

		res, _ := testEvalProgram(tr, tt.input)

		tr.AssertNotNil(res)
		tr.AssertEqual(res.Type(), tt.expectedType)

		switch n := res.(type) {
		case *object.Boolean:
			tr.AssertEqual(n.Value, tt.expectedValue)
		case *object.Number:
			tr.AssertEqual(n.Value, tt.expectedValue)
		case *object.Error:
			// TODO: find a way to test errors
			break

		default:
			tr.T.Errorf("[%d] unexpected type %T", idx, n)
		}
	}
}

func TestLetStatement(t *testing.T) {

	input := "let ident = 5 + 5"

	tr := tester.New(t, input)

	res, e := testEvalProgram(tr, input)

	value := e.getVar("ident")

	tr.AssertEqual(res.Type(), object.NUMBER_OBJ)
	tr.AssertNotNil(value)
	tr.AssertEqual(value.Type(), object.NUMBER_OBJ)
	tr.AssertEqual(value.(*object.Number).Value, float64(10))
}

func TestBlockStatement(t *testing.T) {
	input := `
	let a = 10
	{
		let a = 5
	}
	`
	tr := tester.New(t, "")

	inner, env := testEvalProgram(tr, input)
	tr.AssertNotNil(env)

	outer := env.getVar("a")

	tr.SetName("testing outer")
	tr.AssertNotNil(outer)
	tr.AssertEqual(outer.Type(), object.NUMBER_OBJ)
	tr.AssertEqual(outer.(*object.Number).Value, float64(10))

	tr.SetName("testing inner")
	tr.AssertNotNil(inner)
	tr.AssertEqual(inner.Type(), object.NUMBER_OBJ)
	tr.AssertEqual(inner.(*object.Number).Value, float64(5))

}

func TestIfStatement(t *testing.T) {
	input := `
	let a = 10
	if a > 5 {
		true
	} else {
		false
	}
	`
	tr := tester.New(t, "")

	res, _ := testEvalProgram(tr, input)

	tr.AssertNotNil(res)
	tr.AssertEqual(res.Type(), object.BOOL_OBJ)
	tr.AssertEqual(res.(*object.Boolean).Value, true)

}
func TestElseStatement(t *testing.T) {
	input := `
	let a = 10
	if a < 5 {
		true
	} else {
		false
	}
	`
	tr := tester.New(t, "")

	res, _ := testEvalProgram(tr, input)

	tr.AssertNotNil(res)
	tr.AssertEqual(res.Type(), object.BOOL_OBJ)
	tr.AssertEqual(res.(*object.Boolean).Value, false)
}

func TestIdentifer(t *testing.T) {

	input := `
	let a = 10
	let b = 20
	a + b
	`

	tr := tester.New(t, input)

	res, e := testEvalProgram(tr, input)

	a := e.getVar("a")
	b := e.getVar("b")
	tr.AssertNotNil(b)

	tr.SetName("testing value `a`")
	tr.AssertNotNil(a)
	tr.AssertEqual(a.Type(), object.NUMBER_OBJ)
	tr.AssertEqual(a.(*object.Number).Value, float64(10))

	tr.SetName("testing value `b`")
	tr.AssertNotNil(b)
	tr.AssertEqual(b.Type(), object.NUMBER_OBJ)
	tr.AssertEqual(b.(*object.Number).Value, float64(20))

	tr.SetName(`testing result`)
	tr.AssertNotEqual(res, NIL)
	tr.AssertEqual(res.Type(), object.NUMBER_OBJ)
	tr.AssertEqual(res.(*object.Number).Value, float64(30))

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

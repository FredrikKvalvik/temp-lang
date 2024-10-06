package interpreter

import (
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

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {

			result, _ := testEvalProgram(tr, tt.input)

			tr.AssertEqual(result.Type(), tt.expectedType)

			testAssertType(tr, result, tt.expectedType, tt.expectedVal)
		})
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

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			res, _ := testEvalProgram(tr, tt.input)

			tr.AssertNotNil(res)
			tr.AssertEqual(res.Type(), tt.expectedType)

			testAssertType(tr, res, tt.expectedType, tt.expectedValue)

		})
	}
}

func TestLetStatement(t *testing.T) {

	input := "let ident = 5 + 5"

	tr := tester.New(t, input)

	res, e := testEvalProgram(tr, input)

	value := e.GetVar("ident")

	tr.AssertEqual(res.Type(), object.NUMBER_OBJ)
	tr.AssertNotNil(value)
	tr.AssertEqual(value.Type(), object.NUMBER_OBJ)
	tr.AssertEqual(value.(*object.Number).Value, float64(10))
}

func TestAssignment(t *testing.T) {

	tests := []struct {
		input         string
		varName       string
		expectedType  object.ObjectType
		expectedValue any
	}{
		{"let a = 10; a = 100",
			"a",
			object.NUMBER_OBJ, float64(100),
		},
		{`let b = 10; b = "hello"`,
			"b",
			object.STRING_OBJ, "hello",
		},
		{`let c = ""
			{
				c = "from scope"
			}`,
			"c",
			object.STRING_OBJ, "from scope",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tr := tester.New(t, "")

			res, e := testEvalProgram(tr, tt.input)
			if res.Type() == object.ERROR_OBJ {
				tr.T.Log(res.Inspect())
			}

			value := e.GetVar(tt.varName)
			tr.AssertNotNil(value)
			tr.AssertEqual(value.Type(), tt.expectedType)

			testAssertType(tr, value, tt.expectedType, tt.expectedValue)

		})
	}
}

func TestFnLiteralExpression(t *testing.T) {

	input := "let a = fn() { 10; }"

	tr := tester.New(t, input)

	res, e := testEvalProgram(tr, input)

	value := e.GetVar("a")

	tr.AssertEqual(res.Type(), object.FUNCTION_LITERAL_OBJ)
	tr.AssertNotNil(value)
	tr.AssertEqual(value.Type(), object.FUNCTION_LITERAL_OBJ)
	tr.AssertEqual(len(value.(*object.FnLiteral).Parameters), 0)
	tr.AssertEqual(len(value.(*object.FnLiteral).Body.Statements), 1)
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

	outer := env.GetVar("a")

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

	a := e.GetVar("a")
	b := e.GetVar("b")
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

func testAssertType(
	tr *tester.Tester,
	value object.Object,
	expectedType object.ObjectType,
	expectedValue any,
) {
	tr.T.Helper()

	switch expectedType {
	case object.NUMBER_OBJ:
		tr.AssertEqual(value.(*object.Number).Value, expectedValue)
	case object.STRING_OBJ:
		tr.AssertEqual(value.(*object.String).Value, expectedValue)
	case object.BOOL_OBJ:
		tr.AssertEqual(value.(*object.Boolean).Value, expectedValue)
	case object.ERROR_OBJ:
		// TODO: check error message
	default:
		tr.T.Errorf("Unexpected object type=%s\n", value.Type())
	}
}

func testEvalProgram(tr *tester.Tester, input string) (object.Object, *object.Environment) {
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

	env := object.NewEnv(nil)

	return Eval(program, env), env
}

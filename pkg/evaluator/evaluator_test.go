package evaluator

import (
	"errors"
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

		{`10 + 2 * "100"`,
			nil, object.ERROR_OBJ},

		// boolean returns
		{"10 == 2",
			false, object.BOOL_OBJ},
		{"10 != 2",
			true, object.BOOL_OBJ},
		{`10 != "hello"`,
			true, object.BOOL_OBJ},
		{`10 == "hello"`,
			false, object.BOOL_OBJ},
		{`"hello" == "hello"`,
			true, object.BOOL_OBJ},
		{`"hello" != "goodbye"`,
			true, object.BOOL_OBJ},
		{"10 < 2",
			false, object.BOOL_OBJ},
		{"10 > 2",
			true, object.BOOL_OBJ},
		{`10 > "5"`,
			nil, object.ERROR_OBJ},

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

func TestLogicalExpression(t *testing.T) {

	tests := []struct {
		input       string
		expectedVal bool
	}{
		{"true and false",
			false},
		{"true or false",
			true},
		{"true and true",
			true},
		{"true and true and true and true and true and false",
			false},
		{"10 > 5 and 2 < 4",
			true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tr := tester.New(t, "")

			result, _ := testEvalProgram(tr, tt.input)
			tr.T.Log(result)

			testAssertType(tr, result, object.BOOL_OBJ, tt.expectedVal)
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

	value := e.FindVar("ident")

	tr.AssertEqual(res.Type(), object.NUMBER_OBJ)
	tr.AssertNotNil(value)
	tr.AssertEqual(value.Type(), object.NUMBER_OBJ)
	tr.AssertEqual(value.(*object.NumberObj).Value, float64(10))
}

func TestAssignment(t *testing.T) {

	tests := []struct {
		input         string
		expectedType  object.ObjectType
		expectedValue any
	}{
		{"let a = 10; a = 100",
			object.NUMBER_OBJ, float64(100),
		},
		{`let b = 10; b = "hello"`,
			object.STRING_OBJ, "hello",
		},
		{`let c = ""
			{
				c = "from scope"
			}`,
			object.STRING_OBJ, "from scope",
		},
		{`let c = ["outer"]
			{
				c[0] = "list value assigned from scope"
			}
			 c[0]`,
			object.STRING_OBJ, "list value assigned from scope",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tr := tester.New(t, "")

			res, _ := testEvalProgram(tr, tt.input)
			fmt.Printf("res: %v\n", res)
			if res.Type() == object.ERROR_OBJ {
				tr.T.Log(res.Inspect())
			}

			tr.AssertNotNil(res)
			tr.AssertEqual(res.Type(), tt.expectedType)

			testAssertType(tr, res, tt.expectedType, tt.expectedValue)

		})
	}
}

func TestFnLiteralExpression(t *testing.T) {

	input := "let a = fn() { 10; }"

	tr := tester.New(t, input)

	res, e := testEvalProgram(tr, input)

	value := e.FindVar("a")

	tr.AssertEqual(res.Type(), object.FUNCTION_LITERAL_OBJ)
	tr.AssertNotNil(value)
	tr.AssertEqual(value.Type(), object.FUNCTION_LITERAL_OBJ)
	tr.AssertEqual(len(value.(*object.FnLiteralObj).Parameters), 0)
	tr.AssertEqual(len(value.(*object.FnLiteralObj).Body.Statements), 1)
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

	outer := env.FindVar("a")

	tr.SetName("testing outer")
	tr.AssertNotNil(outer)
	tr.AssertEqual(outer.Type(), object.NUMBER_OBJ)
	tr.AssertEqual(outer.(*object.NumberObj).Value, float64(10))

	tr.SetName("testing inner")
	tr.AssertNotNil(inner)
	tr.AssertEqual(inner.Type(), object.NUMBER_OBJ)
	tr.AssertEqual(inner.(*object.NumberObj).Value, float64(5))

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
	tr.AssertEqual(res.(*object.BooleanObj).Value, true)

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
	tr.AssertEqual(res.(*object.BooleanObj).Value, false)
}

func TestIdentifer(t *testing.T) {

	input := `
	let a = 10
	let b = 20
	a + b
	`

	tr := tester.New(t, input)

	res, e := testEvalProgram(tr, input)

	a := e.FindVar("a")
	b := e.FindVar("b")
	tr.AssertNotNil(b)

	tr.SetName("testing value `a`")
	tr.AssertNotNil(a)
	tr.AssertEqual(a.Type(), object.NUMBER_OBJ)
	tr.AssertEqual(a.(*object.NumberObj).Value, float64(10))

	tr.SetName("testing value `b`")
	tr.AssertNotNil(b)
	tr.AssertEqual(b.Type(), object.NUMBER_OBJ)
	tr.AssertEqual(b.(*object.NumberObj).Value, float64(20))

	tr.SetName(`testing result`)
	tr.AssertNotEqual(res, NIL)
	tr.AssertEqual(res.Type(), object.NUMBER_OBJ)
	tr.AssertEqual(res.(*object.NumberObj).Value, float64(30))

}

// collection of tests for all builtin functions
func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{`len("")`, float64(0)},
		{`len("hello")`, float64(5)},
		{`len("hello", "world")`, object.ArityError},
		{`len()`, object.ArityError},
		{`len([])`, float64(0)},
		{`len([1,2,3])`, float64(3)},
		{`len({})`, float64(0)},
		{`len({true: false, 1: 2})`, float64(2)},

		{`push([], 1)`, []float64{1}},
		{`push([], 1, 2, 3)`, []float64{1, 2, 3}},
		{`push([])`, object.ArityError},
		{`push()`, object.ArityError},
		{`push([1, 2], 3)`, []float64{1, 2, 3}},

		{`pop([1])`, float64(1)},
		{`pop([2, 1])`, float64(1)},
		{`pop([])`, NIL},
		{`pop({"in": "valid"})`, object.TypeError},

		{`str({})`,
			`{
}`},
		{`str([])`, `[]`},
		{`str(10)`, `10`},
		{`str([1,2,3])`, `[1, 2, 3]`},
		{`str("hello world")`, `"hello world"`},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tr := tester.New(t, "")
			result, _ := testEvalProgram(tr, tt.input)

			tr.T.Log(result.Inspect())

			switch tt.expected.(type) {
			case float64:
				tr.AssertEqual(result.Type(), object.NUMBER_OBJ, "result type must equal NUMBER_OBJ")
				tr.AssertEqual(result.(*object.NumberObj).Value, tt.expected, "result must equal expected value")

			case string:
				tr.AssertEqual(result.Type(), object.STRING_OBJ, "result type must equal STRING_OBJ")
				tr.AssertEqual(result.(*object.StringObj).Value, tt.expected, "result must equal expected value")

			case []float64:
				tr.AssertEqual(result.Type(), object.LIST_OBJ)
				values := result.(*object.ListObj).Values
				expected := tt.expected.([]float64)

				tr.AssertEqual(len(values), len(expected))
				for idx, eVal := range expected {
					testAssertType(tr, values[idx], object.NUMBER_OBJ, eVal)
				}

			case *object.NilObj:
				tr.AssertEqual(result.Type(), object.NIL_OBJ)
				tr.AssertEqual(result, tt.expected)

			case error:
				tr.AssertEqual(result.Type(), object.ERROR_OBJ, "result type must equal expected type")
				err := result.(*object.ErrorObj).Error
				tr.AssertTrue(errors.Is(err, tt.expected.(error)), "assert that error is of correct type")

			default:
				tr.T.Fatalf("uncovered test case for type: %T", tt.expected)
			}
		})
	}
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
		tr.AssertEqual(value.(*object.NumberObj).Value, expectedValue)
	case object.STRING_OBJ:
		tr.AssertEqual(value.(*object.StringObj).Value, expectedValue)
	case object.BOOL_OBJ:
		tr.AssertEqual(value.(*object.BooleanObj).Value, expectedValue)
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

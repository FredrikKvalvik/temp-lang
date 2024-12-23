package parser

import (
	"errors"
	"fmt"
	"testing"

	"github.com/fredrikkvalvik/temp-lang/pkg/ast"
	"github.com/fredrikkvalvik/temp-lang/pkg/lexer"
	"github.com/fredrikkvalvik/temp-lang/pkg/tester"
	"github.com/fredrikkvalvik/temp-lang/pkg/token"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input         string
		exptectedName string
		expectedValue any
	}{
		{"let a = 1", "a", float64(1)},
		{"let a = 1.25", "a", float64(1.25)},
		{`let a = ""`, "a", ""},
		{"let ident = true", "ident", true},
		{`let fre = "hei"`, "fre", "hei"},
		{`let      hei_ha         = "hei"`, "hei_ha", "hei"},
		{`let ___hei = "hei"`, "___hei", "hei"},
	}

	for i, tt := range tests {
		p := testParseProgram(tt.input)

		if len(p.Statements) != 1 {
			t.Fatalf("[t: %d] wrong number of statements, got=%d, expected=%d\n", i, len(p.Statements), 1)
		}
		stmt := p.Statements[0]

		if !testLetStatement(t, i, stmt, tt.exptectedName) {
			return
		}

		val := stmt.(*ast.LetStmt).Value
		if !testLiteralExpression(t, i, val, tt.expectedValue) {
			return
		}
	}
}
func testLetStatement(t *testing.T, i int, stmt ast.Stmt, expectedName string) bool {
	let, ok := stmt.(*ast.LetStmt)
	if !ok {
		t.Fatalf("[t: %d] could not assert stmt as *ast.LetStmt\n", i)
		return false
	}

	if let.Name.Value != expectedName {
		t.Errorf("[t: %d] expected name='%+s', got='%+s'\n", i, expectedName, let.Name.Value)
		return false
	}

	if let.Name.Lexeme() != expectedName {
		t.Errorf("[t: %d] expected lexeme='%s', got='%s'\n", i, expectedName, let.Name.Lexeme())
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue any
	}{
		{`return "hello"`, "hello"},
		{`return 100`, float64(100)},
		{`return`, nil},
	}

	for i, tt := range tests {
		tr := tester.New(t, fmt.Sprintf("%d", i))
		p := testParseProgram(tt.input)

		tr.AssertEqual(len(p.Statements), 1)
		retStmt, ok := p.Statements[0].(*ast.ReturnStmt)

		tr.AssertEqual(ok, true)

		switch ev := tt.expectedValue.(type) {
		case string:
		case float64:
			tr.AssertEqual(testLiteralExpression(t, i, retStmt.Value, ev), true)
		default:
			tr.AssertNil(ev)
		}
	}
}

// TODO: finish implementing test
func TestIfStatement(t *testing.T) {
	input := `if x { y; }`

	p := testParseProgram(input)

	if len(p.Statements) != 1 {
		t.Fatalf("exptected len(p.Statements)=1, got=%d\n", len(p.Statements))
	}

	ifStmt, ok := p.Statements[0].(*ast.IfStmt)
	if !ok {
		t.Fatalf("could not assert stmt type=ast.IfStmt\n")
	}

	x, ok := ifStmt.Condition.(*ast.IdentifierExpr)
	if !ok {
		t.Fatalf("could not assert condition type=ast.IdentifierExpr\n")
	}

	if x.Value != "x" {
		t.Fatalf("x.value expected='x', got='%s'", x.Value)
	}

	if len(ifStmt.Then.Statements) != 1 {
		t.Fatalf("len(IfStmt.Then.Statements) expected=1, got=%d\n", len(ifStmt.Then.Statements))
	}

	exprStmt, ok := ifStmt.Then.Statements[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("could not assert condition type=ast.ExpressionStatment\n")
	}

	y, ok := exprStmt.Expression.(*ast.IdentifierExpr)
	if !ok {
		t.Fatalf("could not assert condition type=ast.IdentifierExpr\n")
	}

	if y.Value != "y" {
		t.Fatalf("x.value expected='x', got='%s'", y.Value)
	}
}

func TestExpressions(t *testing.T) {

	tests := []struct {
		input    string
		expected string
	}{
		{"1+2+3",
			"((1 + 2) + 3)"},
		{"1+2*3",
			"(1 + (2 * 3))"},
		{"true and false or true",
			"((true and false) or true)"},
		{"true or false and true",
			"(true or (false and true))"},
		{"true and false and true",
			"((true and false) and true)"},
		{"1 < 2 > 3",
			"((1 < 2) > 3)"},
		{"(1 + 2) * 3",
			"((1 + 2) * 3)"},
		{"1 * (2 + 3)",
			"(1 * (2 + 3))"},
	}

	for i, tt := range tests {
		p := testParseProgram(tt.input)

		if len(p.Statements) != 1 {
			t.Fatalf("[t: %d] expected p.Statments len=1, got=%d\n", i, len(p.Statements))
		}

		stmt, ok := p.Statements[0].(*ast.ExpressionStmt)
		if !ok {
			t.Fatalf("[t: %d] could not assert type *ast.ExpressionStmt\n", i)
		}

		str := stmt.Expression.String()
		if str != tt.expected {
			t.Errorf("[t: %d] expected=%s, got=%s\n", i, tt.expected, str)
		}
	}

}
func TestBinaryExpression(t *testing.T) {

	tests := []struct {
		input string
		left  any
		op    token.TokenType
		right any
	}{
		{"5 + 5", float64(5), token.PLUS, float64(5)},
		{"5 - 5", float64(5), token.MINUS, float64(5)},
		{"5 * 5", float64(5), token.ASTERISK, float64(5)},
		{"5 / 5", float64(5), token.SLASH, float64(5)},
		// TODO: add test for logicalParsing
		// {"5 and 5", float64(5), token.AND, float64(5)},
		// {"5 or 5", float64(5), token.OR, float64(5)},
		{"5 > 5", float64(5), token.GT, float64(5)},
		{"5 < 5", float64(5), token.LT, float64(5)},
		{"5 == 5", float64(5), token.EQ, float64(5)},
		{"5 != 5", float64(5), token.NOT_EQ, float64(5)},
	}

	for i, tt := range tests {
		p := testParseProgram(tt.input)

		if len(p.Statements) != 1 {
			t.Fatalf("[t: %d] expected p.Statments len=1, got=%d\n", i, len(p.Statements))
		}

		expr, ok := p.Statements[0].(*ast.ExpressionStmt).Expression.(*ast.BinaryExpr)
		if !ok {
			t.Fatalf("[t: %d] could not assert type *ast.BinaryExpr\n", i)
		}

		if !testBinaryExpression(t, i, expr, tt.left, tt.op, tt.right) {
			return
		}
	}

}

func TestWhileStatement(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"while {}"},
		{"while true {}"},
		{"while false {}"},
		{"while 1 > 0 {}"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tr := tester.New(t, "")
			res := testParseProgram(tt.input)

			tr.AssertNotNil(res, "the program must be non-nil")
			tr.AssertEqual(len(res.Statements), 1, "expect a single stmt")
			whileStmt, ok := res.Statements[0].(*ast.WhileStmt)
			tr.AssertTrue(ok, "statment must be whileStmt")

			tr.AssertNotNil(whileStmt.Body)
		})
	}
}

func TestFunctionLiterals(t *testing.T) {

	tests := []struct {
		input           string
		expectedArgLen  int
		expectedBodyLen int
	}{
		{"(fn() {})",
			0, 0},
		{"(fn() { 10; })",
			0, 1},
		{"(fn(a) {})",
			1, 0},
		{"(fn(a, b, c) { })",
			3, 0},
		{"(fn(a, b, c) { 1; 2; 3})",
			3, 3},
	}

	for idx, tt := range tests {
		tr := tester.New(t, fmt.Sprintf("[%d]", idx))

		res := testParseProgram(tt.input)

		tr.AssertNotNil(res)
		tr.AssertEqual(len(res.Statements), 1)

		tr.SetName(fmt.Sprintf("[%d]is expression", idx))
		expr, ok := res.Statements[0].(*ast.ExpressionStmt)
		tr.AssertEqual(ok, true)
		tr.AssertNotNil(expr)

		tr.SetName(fmt.Sprintf("[%d]is functionLiteral", idx))
		fun, ok := expr.Expression.(*ast.ParenExpr).Expression.(*ast.FunctionLiteralExpr)
		tr.AssertTrue(ok)

		tr.SetName(fmt.Sprintf("[%d]test function args", idx))
		tr.AssertEqual(len(fun.Arguments), tt.expectedArgLen)

		tr.SetName(fmt.Sprintf("[%d]test function body", idx))
		tr.AssertNotNil(fun.Body)
		tr.AssertEqual(len(fun.Body.Statements), tt.expectedBodyLen)
	}
}
func TestFunctionCalls(t *testing.T) {

	// TODO: can probalby improve test by testing the arguments of
	// the function calls, instead of just the number of args
	tests := []struct {
		input           string
		expectedCallee  string
		expectedArgsLen int
	}{
		// without args
		{"function()",
			"function", 0},
		{"function(1)",
			"function", 1},
		// with args
		{"function(a,b,c)",
			"function", 3},
		// anynomous function call without args
		{"(fn() {})()",
			"fn() {\n}", 0},
		// anynomous function call without args
		{"(fn(a) {})(1)",
			"fn(a) {\n}", 1},
		// anynomous function call with args
		{"(fn(a,b,c) {})(1,2,3)",
			"fn(a, b, c) {\n}", 3},
	}

	for idx, tt := range tests {
		tr := tester.New(t, fmt.Sprintf("[%d]", idx))

		res := testParseProgram(tt.input)

		tr.AssertNotNil(res)
		tr.AssertEqual(len(res.Statements), 1)

		tr.SetName(fmt.Sprintf("[%d]is expression", idx))
		expr, ok := res.Statements[0].(*ast.ExpressionStmt)
		tr.AssertEqual(ok, true)
		tr.AssertNotNil(expr)

		tr.SetName(fmt.Sprintf("[%d]is call", idx))
		call, ok := expr.Expression.(*ast.CallExpr)
		tr.AssertTrue(ok)

		tr.SetName(fmt.Sprintf("[%d]test callee name", idx))
		tr.AssertEqual(call.Callee.String(), tt.expectedCallee)

		tr.SetName(fmt.Sprintf("[%d]test call args", idx))
		tr.AssertEqual(len(call.Arguments), tt.expectedArgsLen)

	}
}

func TestListLiteralExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected []any
	}{
		{"[]",
			[]any{}},
		{"[1]",
			[]any{float64(1)}},
		{`[1, "hello", "world"]`,
			[]any{float64(1), "hello", "world"}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tr := tester.New(t, "")

			res := testParseProgram(tt.input)
			tr.AssertEqual(len(res.Statements), 1)

			expression, ok := res.Statements[0].(*ast.ExpressionStmt)
			tr.AssertTrue(ok)

			listLit, ok := expression.Expression.(*ast.ListLiteralExpr)
			tr.AssertTrue(ok)

			tr.AssertEqual(len(listLit.Items), len(tt.expected))
		})
	}
}

func TestMapLiteralExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected map[any]any
	}{
		{"({})",
			map[any]any{}},
		{`({"10":10})`,
			map[any]any{"10": 10}},
		{`({1:10,true:false})`,
			map[any]any{float64(1): 10, true: false}},
		{`({
			1:10,
			true:false,
			})`,
			map[any]any{float64(1): 10, true: false}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tr := tester.New(t, "")

			res := testParseProgram(tt.input)
			tr.AssertEqual(len(res.Statements), 1)

			expression, ok := res.Statements[0].(*ast.ExpressionStmt)
			tr.AssertTrue(ok)

			mapLit, ok := expression.Expression.(*ast.ParenExpr).Expression.(*ast.MapLiteralExpr)
			tr.AssertTrue(ok)

			tr.AssertEqual(len(mapLit.KeyValues), len(tt.expected))
		})
	}
}

func TestImportStmt(t *testing.T) {
	tests := []struct {
		input        string
		expectedName string
		expectedPath string
	}{
		{`import ident "path"`,
			"ident", "path"},
		{`import http "http";`,
			"http", "http"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tr := tester.New(t, "")
			l := lexer.New(tt.input)
			p := New(l)

			res := p.ParseProgram()

			if len(p.Errors()) > 0 {
				tr.T.Log("did log")
				tr.T.Errorf("parser error\n%s", errors.Join(p.errors...))
			}

			tr.AssertEqual(len(res.Statements), 1, "expect 1 stamtement in program")

			imp := res.Statements[0].(*ast.ImportStmt)
			tr.AssertEqual(imp.GetToken().Type, token.IMPORT, "expect the statement to be ImportStmt")

			tr.AssertEqual(imp.Name.Value, tt.expectedName, fmt.Sprintf("expect import name: %s == %s", imp.Name.Value, tt.expectedName))
			tr.AssertEqual(imp.Path, tt.expectedPath, fmt.Sprintf("expect import path: %s == %s", imp.Path, tt.expectedPath))
		})
	}
}

func literalToValue(expr ast.Expr) any {
	switch e := expr.(type) {
	case *ast.NumberLiteralExpr:
		return e.Value
	case *ast.StringLiteralExpr:
		return e.Value
	case *ast.BooleanLiteralExpr:
		return e.Value
	}

	return nil
}

func TestIndexExpressions(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"a[0]"},
		{"[list][0]"},
		{`"hello world"[0]`},
		{`({"hello": "world"})["hello"]`},
		{`({"hello": "world"})[true]`},
		{`({"hello": "world"})[false]`},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tr := tester.New(t, "")

			res := testParseProgram(tt.input)
			tr.AssertEqual(len(res.Statements), 1)

			expression, ok := res.Statements[0].(*ast.ExpressionStmt)
			tr.AssertTrue(ok)

			indexExpr, ok := expression.Expression.(*ast.IndexExpr)
			tr.AssertTrue(ok)

			tr.AssertNotNil(indexExpr.Left)
			tr.AssertNotNil(indexExpr.Index)
		})
	}
}

func testBinaryExpression(t *testing.T, i int, expr *ast.BinaryExpr, eLeft any, op token.TokenType, eRight any) bool {
	t.Helper()
	if !testLiteralExpression(t, i, expr.Left, eLeft) {
		t.Errorf("[t: %d] Left: expected=%v, got=%v\n", i, eLeft, expr.Left.Lexeme())
		return false
	}
	if expr.Operand != op {
		t.Errorf("[t: %d] Op: expected=%v, got=%v\n", i, op, expr.Operand.String())
		return false
	}

	if !testLiteralExpression(t, i, expr.Right, eRight) {
		t.Errorf("[t: %d] Right: expected=%v, got=%v\n", i, eRight, expr.Right.Lexeme())
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, i int, expr ast.Expr, expectedValue any) bool {
	t.Helper()

	switch val := expectedValue.(type) {
	case float64:
		return testNumberLiteral(t, i, expr, val)
	case string:
		return testStringLiteral(t, i, expr, val)
	case bool:
		return testBooleanLiteral(t, i, expr, val)
	default:
		t.Errorf("[t: %d] unexpected expr type=%s", i, expr.Lexeme())
		return false
	}
}

func testStringLiteral(t *testing.T, i int, expr ast.Expr, expectedValue string) bool {
	t.Helper()
	strLit, ok := expr.(*ast.StringLiteralExpr)
	if !ok {
		t.Errorf("[t: %d] expr is not *ast.StringLiteralExpr\n", i)
		return false
	}
	if strLit.Value != expectedValue {
		t.Errorf("[t: %d] expected value=%s, got=%s\n", i, expectedValue, strLit.Value)
		return false
	}

	return true
}
func testNumberLiteral(t *testing.T, i int, expr ast.Expr, expectedValue float64) bool {
	t.Helper()
	numLit, ok := expr.(*ast.NumberLiteralExpr)
	if !ok {
		t.Errorf("[t: %d] expr is not *ast.NumberLiteralExpr\n", i)
		return false
	}
	if numLit.Value != expectedValue {
		t.Errorf("[t: %d] expected value=%f, got=%f\n", i, expectedValue, numLit.Value)
		return false
	}

	return true
}
func testBooleanLiteral(t *testing.T, i int, expr ast.Expr, expectedValue bool) bool {
	t.Helper()
	boolLit, ok := expr.(*ast.BooleanLiteralExpr)
	if !ok {
		t.Errorf("[t: %d] expr is not *ast.BooleanLiteralExpr\n", i)
		return false
	}
	if boolLit.Value != expectedValue {
		t.Errorf("[t: %d] expected value=%t, got=%t\n", i, expectedValue, boolLit.Value)
		return false
	}

	return true
}

func testParseProgram(input string) *ast.Program {
	l := lexer.New(input)
	p := New(l)

	return p.ParseProgram()
}

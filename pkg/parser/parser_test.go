package parser

import (
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
		{"5 and 5", float64(5), token.AND, float64(5)},
		{"5 or 5", float64(5), token.OR, float64(5)},
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

func TestFunctionLiteral(t *testing.T) {
	tr := tester.New(t, "")

	input := "fn() { 10; }"

	res := testParseProgram(input)

	tr.AssertNotNil(res)
	tr.AssertEqual(len(res.Statements), 1)

	tr.SetName("is expression")
	expr, ok := res.Statements[0].(*ast.ExpressionStmt)
	tr.AssertEqual(ok, true)
	tr.AssertNotNil(expr)

	tr.SetName("is functionLiteral")
	fun, ok := expr.Expression.(*ast.FunctionLiteralExpr)
	tr.AssertTrue(ok)

	tr.SetName("test function args")
	tr.AssertEqual(len(fun.Arguments), 0)

	tr.SetName("test function body")
	tr.AssertNotNil(fun.Body)
	tr.AssertEqual(len(fun.Body.Statements), 1)

	expr, ok = fun.Body.Statements[0].(*ast.ExpressionStmt)
	tr.AssertTrue(ok)
	tr.AssertNotNil(expr)
	num, ok := expr.Expression.(*ast.NumberLiteralExpr)
	tr.AssertTrue(ok)
	tr.AssertEqual(num.Value, float64(10))
}

func testBinaryExpression(t *testing.T, i int, expr *ast.BinaryExpr, eLeft any, op token.TokenType, eRight any) bool {
	t.Helper()
	if !testLiteralExpression(t, i, expr.Left, eLeft) {
		t.Errorf("[t: %d] Left: expected=%v, got=%v\n", i, eLeft, expr.Left.Literal())
		return false
	}
	if expr.Operand != op {
		t.Errorf("[t: %d] Op: expected=%v, got=%v\n", i, op, expr.Operand.String())
		return false
	}

	if !testLiteralExpression(t, i, expr.Right, eRight) {
		t.Errorf("[t: %d] Right: expected=%v, got=%v\n", i, eRight, expr.Right.Literal())
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

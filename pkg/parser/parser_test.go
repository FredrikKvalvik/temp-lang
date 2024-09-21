package parser

import (
	"testing"

	"github.com/fredrikkvalvik/temp-lang/pkg/ast"
	"github.com/fredrikkvalvik/temp-lang/pkg/lexer"
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
		{`let hei_ha = "hei"`, "hei_ha", "hei"},
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

func TestExpressionStickiness(t *testing.T) {

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

func testLiteralExpression(t *testing.T, i int, expr ast.Expr, expectedValue any) bool {

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

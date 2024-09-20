package parser

import (
	"testing"

	"github.com/fredrikkvalvik/temp-lang/pkg/ast"
	"github.com/fredrikkvalvik/temp-lang/pkg/lexer"
)

func TestLetStatement(t *testing.T) {
	input := `
let a = "string"
let b = 5;
let c = 1+2/2
`
	program := testGenerateProgram(t, input)

	if len(program.Statements) != 3 {
		t.Fatalf("program len is not 3, len=%d", len(program.Statements))
	}

}

func testGenerateProgram(t *testing.T, input string) *ast.Program {
	l := lexer.New(input)

	if l.DidError() {
		for _, err := range l.Errors() {
			t.Error(err.Error())
		}
		t.Fatalf("lexer could not parse input")
	}
	p := New(l)

	program := p.ParseProgram()

	return program
}

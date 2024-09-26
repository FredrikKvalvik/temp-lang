package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/fredrikkvalvik/temp-lang/pkg/interpreter"
	"github.com/fredrikkvalvik/temp-lang/pkg/lexer"
	"github.com/fredrikkvalvik/temp-lang/pkg/parser"
)

func main() {
	args := os.Args[1:]

	if len(args) > 1 {
		fmt.Printf("usage: templang [script-path]")
		os.Exit(1)
	} else if len(args) == 1 {
		fmt.Printf("running scripts: %s...", args[0])
	} else {
		repl(os.Stdin, os.Stdout)
	}
}

func repl(in io.Reader, out io.Writer) {
	s := bufio.NewScanner(in)

	env := interpreter.NewEnv(nil)
	for {
		fmt.Print("> ")
		scanned := s.Scan()
		if !scanned {
			return
		}

		line := s.Text()
		l := lexer.New(line)
		if l.DidError() {
			for _, err := range l.Errors() {
				fmt.Fprintf(out, "%s\n", err)
			}
			continue
		}

		p := parser.New(l)
		program := p.ParseProgram()

		if p.DidError() {
			for _, err := range p.Errors() {
				fmt.Println(err.Error())
			}
			continue
		}

		result := interpreter.Eval(program, env)
		fmt.Printf("%s\n", result.Inspect())
	}
}

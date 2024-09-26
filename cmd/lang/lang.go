package main

import (
	"fmt"
	"os"

	"github.com/fredrikkvalvik/temp-lang/pkg/interpreter"
	"github.com/fredrikkvalvik/temp-lang/pkg/lexer"
	"github.com/fredrikkvalvik/temp-lang/pkg/object"
	"github.com/fredrikkvalvik/temp-lang/pkg/parser"
	"github.com/fredrikkvalvik/temp-lang/pkg/repl"
)

func main() {
	args := os.Args[1:]

	if len(args) > 1 {
		fmt.Printf("usage: templang [script-path]")
		os.Exit(1)
	} else if len(args) == 1 {
		path := args[0]
		file := readFile(path)
		_ = runProgram(file)

	} else {
		env := interpreter.NewEnv(nil)
		repl.New(env).Run(os.Stdin, os.Stdout)
		return
	}
}

func runProgram(in string) object.Object {

	l := lexer.New(in)
	if l.DidError() {
		for _, err := range l.Errors() {
			fmt.Printf("%s\n", err)
		}
		return nil
	}

	p := parser.New(l)
	program := p.ParseProgram()

	if p.DidError() {
		for _, err := range p.Errors() {
			fmt.Println(err.Error())
		}
		return nil
	}

	env := interpreter.NewEnv(nil)
	result := interpreter.Eval(program, env)
	return result
}

func readFile(path string) string {
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(file)
}

package main

import (
	"errors"
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
		res, err := runProgram(file)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if res.Type() == object.ERROR_OBJ {
			fmt.Println(res.Inspect())
			return
		}

	} else {
		env := object.NewEnv(nil)
		repl.New(env).Run(os.Stdin, os.Stdout)
		return
	}
}

func runProgram(in string) (object.Object, error) {

	l := lexer.New(in)
	if l.DidError() {
		errs := ""
		for _, err := range l.Errors() {
			errs += fmt.Sprintf("%s\n", err)
		}
		return nil, errors.New(errs)
	}

	p := parser.New(l)
	program := p.ParseProgram()

	if p.DidError() {
		errs := ""
		for _, err := range p.Errors() {
			errs += fmt.Sprintf("%s\n", err)
		}
		return nil, errors.New(errs)
	}

	env := object.NewEnv(nil)
	result := interpreter.Eval(program, env)
	return result, nil
}

func readFile(path string) string {
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(file)
}

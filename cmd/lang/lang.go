package main

import (
	"errors"
	"fmt"
	"os"

	"flag"

	"github.com/fredrikkvalvik/temp-lang/pkg/evaluator"
	"github.com/fredrikkvalvik/temp-lang/pkg/lexer"
	"github.com/fredrikkvalvik/temp-lang/pkg/object"
	"github.com/fredrikkvalvik/temp-lang/pkg/parser"
	"github.com/fredrikkvalvik/temp-lang/pkg/repl"
	"github.com/fredrikkvalvik/temp-lang/pkg/resolver"
)

func main() {
	attach := flag.Bool("attach", false, "attach repl to a program after its execution")

	flag.Parse()

	if len(flag.Args()) > 0 {
		path := flag.Arg(0)
		file := readFile(path)
		env := object.NewEnv(nil)
		res, err := runProgram(file, env)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if res.Type() == object.ERROR_OBJ {
			fmt.Println(res.Inspect())
			return
		}

		if attach != nil && *attach {
			repl.New(os.Stdin, os.Stdout).Run(env)
		}

	} else {
		env := object.NewEnv(nil)
		repl.New(os.Stdin, os.Stdout).Run(env)
		return
	}
}

func runProgram(in string, env *object.Environment) (object.Object, error) {

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
	// fmt.Printf("%s\n", program)
	// return nil, nil

	r := resolver.New(env)
	r.Resolve(program)
	if len(r.Errors) > 0 {
		return nil, errors.Join(r.Errors...)
	}

	result := evaluator.Eval(program, env)
	return result, nil
}

func readFile(path string) string {
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(file)
}

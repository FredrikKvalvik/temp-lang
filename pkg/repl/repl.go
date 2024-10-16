package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/fredrikkvalvik/temp-lang/pkg/evaluator"
	"github.com/fredrikkvalvik/temp-lang/pkg/lexer"
	"github.com/fredrikkvalvik/temp-lang/pkg/object"
	"github.com/fredrikkvalvik/temp-lang/pkg/parser"
	"github.com/fredrikkvalvik/temp-lang/pkg/resolver"
)

type Repl struct {
	// env *object.Environment
	in  io.Reader
	out io.Writer
}

func New(in io.Reader, out io.Writer) *Repl {
	return &Repl{
		in:  in,
		out: out,
	}
}

// os.Stdin and os.Stdout are the usual args
func (r *Repl) Run(env *object.Environment) {
	s := bufio.NewScanner(r.in)

	resolve := resolver.New(env)
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
				fmt.Fprintf(r.out, "%s\n", err)
			}
			continue
		}

		p := parser.New(l)
		program := p.ParseProgram()

		resolve.Resolve(program)
		if len(resolve.Errors) > 0 {
			for _, err := range resolve.Errors {
				fmt.Printf("err: %v\n", err)
			}
			resolve.Errors = []error{}
			continue
		}

		if p.DidError() {
			for _, err := range p.Errors() {
				fmt.Println(err.Error())
			}
			continue
		}

		result := evaluator.Eval(program, env)
		fmt.Printf("%s\n", result.Inspect())
	}
}

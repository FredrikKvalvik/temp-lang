package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/fredrikkvalvik/temp-lang/pkg/interpreter"
	"github.com/fredrikkvalvik/temp-lang/pkg/lexer"
	"github.com/fredrikkvalvik/temp-lang/pkg/object"
	"github.com/fredrikkvalvik/temp-lang/pkg/parser"
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

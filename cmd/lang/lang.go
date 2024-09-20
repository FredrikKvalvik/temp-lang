package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/fredrikkvalvik/temp-lang/pkg/lexer"
	"github.com/fredrikkvalvik/temp-lang/pkg/parser"
	"github.com/fredrikkvalvik/temp-lang/pkg/token"
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
	fmt.Printf("INSIDE REPL\n")
	fmt.Printf("%s\n", token.SEMICOLON)
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

		_, err := fmt.Printf("%s", p.ParseProgram())

		if p.DidError() {
			for _, err := range p.Errors() {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Fprint(out, err)
		}

	}
}

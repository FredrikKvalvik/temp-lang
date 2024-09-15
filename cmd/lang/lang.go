package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/fredrikkvalvik/temp-lang/pkg/scanner"
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
	for {
		fmt.Print("> ")
		scanned := s.Scan()
		if !scanned {
			return
		}
		line := s.Text()
		ts := scanner.New(line)
		var tokens []token.Token
		for {
			tok := ts.NextToken()
			tokens = append(tokens, tok)
			if tok.TokenType == token.EOF {
				break
			}
		}

		fmt.Printf("input: %s\n", line)

		for _, tok := range tokens {
			fmt.Println(tok.String())
		}

		if ts.DidError() {
			for _, err := range ts.Errors() {
				fmt.Println(err)
			}
		}
	}
}

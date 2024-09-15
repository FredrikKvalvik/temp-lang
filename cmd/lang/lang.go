package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
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
	scanner := bufio.NewScanner(in)
	fmt.Printf("INSIDE REPL\n")
	for {
		fmt.Print("> ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Bytes()
		fmt.Fprintf(out, "%s\n", line)
	}
}

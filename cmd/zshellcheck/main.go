package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/lexer"
	"github.com/afadesigns/zshellcheck/pkg/parser"
	"github.com/afadesigns/zshellcheck/pkg/reporter"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: zshellcheck <file1.zsh> [file2.zsh]...")
		os.Exit(1)
	}

	for _, filename := range os.Args[1:] {
		processFile(filename, os.Stdout, os.Stderr)
	}
}

func processFile(filename string, out, errOut io.Writer) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(errOut, "Error reading file %s: %s\n", filename, err)
		return
	}

	l := lexer.New(string(data))
	p := parser.New(l)

	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		for _, msg := range p.Errors() {
			fmt.Fprintf(errOut, "Parser Error in %s: %s\n", filename, msg)
		}
		return
	}

	violations := []katas.Violation{}
	ast.Walk(program, func(node ast.Node) bool {
		violations = append(violations, katas.Check(node)...)
		return true // Continue walking
	})

	if len(violations) > 0 {
		fmt.Fprintf(out, "Violations in %s:\n", filename)
		reporter := reporter.NewTextReporter(out)
		reporter.Report(violations)
	}
}
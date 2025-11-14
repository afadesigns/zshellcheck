package main

import (
	"flag"
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
	format := flag.String("format", "text", "The output format (text or json)")
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("Usage: zshellcheck [flags] <file1.zsh> [file2.zsh]...")
		os.Exit(1)
	}

	for _, filename := range flag.Args() {
		processFile(filename, os.Stdout, os.Stderr, *format)
	}
}

func processFile(filename string, out, errOut io.Writer, format string) {
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
		var r reporter.Reporter
		switch format {
		case "json":
			r = reporter.NewJSONReporter(out)
		default:
			fmt.Fprintf(out, "Violations in %s:\n", filename)
			r = reporter.NewTextReporter(out)
		}
		r.Report(violations)
	}
}

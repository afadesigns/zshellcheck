package main

import (
	"fmt"
	"os"

	"github.com/afadesigns/zshellcheck/pkg/lexer"
	"github.com/afadesigns/zshellcheck/pkg/parser"
	"github.com/afadesigns/zshellcheck/pkg/katas"
)

func main() {
	// For now, we'll hardcode a simple input for testing.
	// In a real CLI, this would come from a file or stdin.
	input := `if [ -f "myfile" ]; then echo "file exists"; fi`

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		for _, msg := range p.Errors() {
			fmt.Fprintf(os.Stderr, "Parser Error: %s\n", msg)
		}
		os.Exit(1)
	}

	fmt.Println("AST:")
	fmt.Println(program.String())
	fmt.Println("\nRunning Katas...")

	// Walk the AST and collect violations
	violations := []katas.Violation{}
	ast.Walk(program, func(node ast.Node) bool {
		for _, kata := range katas.AllKatas {
			v := kata.Check(node)
			violations = append(violations, v...)
		}
		return true // Continue walking
	})

	if len(violations) == 0 {
		fmt.Println("No ZShellCheck Katas violated. Your Zsh is enlightened!")
	} else {
		fmt.Println("ZShellCheck Katas violated:")
		for _, v := range violations {
			fmt.Printf("  %s:%d:%d: [%s] %s\n", "<stdin>", v.Line, v.Column, v.KataID, v.Message)
		}
		os.Exit(1)
	}
}
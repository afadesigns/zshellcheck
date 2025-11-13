package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/lexer"
	"github.com/afadesigns/zshellcheck/pkg/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: zshellcheck <file1.zsh> [file2.zsh]...")
		os.Exit(1)
	}

	for _, filename := range os.Args[1:] {
		processFile(filename)
	}
}

func processFile(filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %s\n", filename, err)
		return
	}

	l := lexer.New(string(data))
	p := parser.New(l)

	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		for _, msg := range p.Errors() {
			fmt.Fprintf(os.Stderr, "Parser Error in %s: %s\n", filename, msg)
		}
		return
	}

	violations := []katas.Violation{}
	ast.Walk(program, func(node ast.Node) bool {
		nodeType := reflect.TypeOf(node)
		if katasForNode, ok := katas.KatasByNodeType[nodeType]; ok {
			for _, kata := range katasForNode {
				v := kata.Check(node)
				violations = append(violations, v...)
			}
		}
		return true // Continue walking
	})

	if len(violations) > 0 {
		fmt.Printf("Violations in %s:\n", filename)
		for _, v := range violations {
			fmt.Printf("  %s:%d:%d: [%s] %s\n", filename, v.Line, v.Column, v.KataID, v.Message)
		}
	}
}
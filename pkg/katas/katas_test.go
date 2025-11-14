package katas

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/lexer"
	"github.com/afadesigns/zshellcheck/pkg/parser"
)

func testFile(t *testing.T, filepath string, kataIDs []string, expectedViolations []Violation) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		t.Fatal(err)
	}

	l := lexer.New(string(data))
	p := parser.New(l)
	program := p.ParseProgram()

	violations := []Violation{}
	ast.Walk(program, func(node ast.Node) bool {
		nodeType := reflect.TypeOf(node)
		if katas, ok := KatasByNodeType[nodeType]; ok {
			for _, kata := range katas {
				for _, id := range kataIDs {
					if kata.ID == id {
						violations = append(violations, kata.Check(node)...)
					}
				}
			}
		}
		return true
	})

	if len(violations) != len(expectedViolations) {
		t.Fatalf("Expected %d violations, got %d", len(expectedViolations), len(violations))
	}

	for i, v := range violations {
		if v.KataID != expectedViolations[i].KataID {
			t.Errorf("Expected KataID %s, got %s", expectedViolations[i].KataID, v.KataID)
		}
		if v.Message != expectedViolations[i].Message {
			t.Errorf("Expected Message %s, got %s", expectedViolations[i].Message, v.Message)
		}
		if v.Line != expectedViolations[i].Line {
			t.Errorf("Expected Line %d, got %d", expectedViolations[i].Line, v.Line)
		}
		if v.Column != expectedViolations[i].Column {
			t.Errorf("Expected Column %d, got %d", expectedViolations[i].Column, v.Column)
		}
	}
}

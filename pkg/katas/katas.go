package katas

import (
	"reflect"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

type Kata struct {
	ID          string
	Title       string
	Description string
	Check       func(node ast.Node) []Violation
}

type Violation struct {
	KataID  string
	Message string
	Line    int
	Column  int
}

var KatasByNodeType = make(map[reflect.Type][]Kata)

func RegisterKata(nodeType reflect.Type, kata Kata) {
	KatasByNodeType[nodeType] = append(KatasByNodeType[nodeType], kata)
}

func Check(node ast.Node) []Violation {
	var violations []Violation
	nodeType := reflect.TypeOf(node)
	if katas, ok := KatasByNodeType[nodeType]; ok {
		for _, kata := range katas {
			violations = append(violations, kata.Check(node)...)
		}
	}
	return violations
}
package katas

import (
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

var AllKatas = []Kata{}

func RegisterKata(kata Kata) {
	AllKatas = append(AllKatas, kata)
}

package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(Kata{
		ID:          "ZC1001",
		Title:       "Prefer [[ over [ for Zsh-specific tests",
		Description: "The [[...]] construct is a Zsh keyword, offering safer and more powerful conditional expressions than the traditional [ command. It prevents word splitting and pathname expansion, and supports advanced features like regex matching.",
		Check:       checkZC1001,
	})
}

func checkZC1001(node ast.Node) []Violation {
	violations := []Violation{}
	return violations
}

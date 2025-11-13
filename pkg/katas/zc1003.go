package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(Kata{
		ID:          "ZC1003",
		Title:       "Prefer [[ over [ for tests",
		Description: "The [[...]] construct is a Zsh keyword, offering safer and more powerful conditional expressions than the traditional [ command. It prevents word splitting and pathname expansion, and supports advanced features like regex matching.",
		Check:       checkZC1003,
	})
}

func checkZC1003(node ast.Node) []Violation {
	violations := []Violation{}

	walkFn := func(node ast.Node) bool {
		if bracketExp, ok := node.(*ast.BracketExpression); ok {
			violations = append(violations, Violation{
				KataID:  "ZC1003",
				Message: "Prefer [[ over [ for tests. [[ is a Zsh keyword that offers safer and more powerful conditional expressions.",
				Line:    bracketExp.Token.Line,
				Column:  bracketExp.Token.Column,
			})
		}
		return true
	}

	ast.Walk(node, walkFn)

	return violations
}

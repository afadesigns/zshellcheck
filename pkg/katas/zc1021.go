package katas

import (

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:          "ZC1021",
		Title:       "Use symbolic permissions with `chmod` instead of octal",
		Description: "Symbolic permissions (e.g., `u+x`) are more readable and less error-prone than " +
			"octal permissions (e.g., `755`).",
		Check:       checkZC1021,
	})
}

func checkZC1021(node ast.Node) []Violation {
	violations := []Violation{}

	if cmd, ok := node.(*ast.SimpleCommand); ok {
		if name, ok := cmd.Name.(*ast.Identifier); ok && name.Value == "chmod" {
			for _, arg := range cmd.Arguments {
				if _, ok := arg.(*ast.IntegerLiteral); ok {
					violations = append(violations, Violation{
						KataID:  "ZC1021",
						Message: "Use symbolic permissions with `chmod` instead of octal.",
						Line:    name.Token.Line,
						Column:  name.Token.Column,
					})
					break
				}
			}
		}
	}

	return violations
}

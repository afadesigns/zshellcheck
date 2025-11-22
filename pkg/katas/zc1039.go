package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:          "ZC1039",
		Title:       "Avoid `rm` with root path",
		Description: "Running `rm` on the root directory `/` is dangerous. " +
			"Ensure you are not deleting the entire filesystem.",
		Check: checkZC1039,
	})
}

func checkZC1039(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	// Check if command is rm
	if cmdName, ok := cmd.Name.(*ast.Identifier); !ok || cmdName.Value != "rm" {
		return nil
	}

	violations := []Violation{}

	for _, arg := range cmd.Arguments {
		if str, ok := arg.(*ast.StringLiteral); ok {
			if str.Value == "/" {
				violations = append(violations, Violation{
					KataID:  "ZC1039",
					Message: "Avoid `rm` on the root directory `/`. This is highly dangerous.",
					Line:    str.Token.Line,
					Column:  str.Token.Column,
				})
			}
		}
	}

	return violations
}

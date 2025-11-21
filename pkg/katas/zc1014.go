package katas

import (

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:          "ZC1014",
		Title:       "Use `git switch` or `git restore` instead of `git checkout`",
		Description: "The `git checkout` command can be ambiguous. `git switch` is used for switching " +
			"branches and `git restore` is used for restoring files. Using these more specific commands " +
			"can make your scripts clearer and less error-prone.",
		Check:       checkZC1014,
	})
}

func checkZC1014(node ast.Node) []Violation {
	violations := []Violation{}

	if cmd, ok := node.(*ast.SimpleCommand); ok {
		if name, ok := cmd.Name.(*ast.Identifier); ok && name.Value == "git" {
			if len(cmd.Arguments) > 0 {
				if arg, ok := cmd.Arguments[0].(*ast.Identifier); ok && arg.Value == "checkout" {
					violations = append(violations, Violation{
						KataID:  "ZC1014",
						Message: "Use `git switch` or `git restore` instead of the ambiguous `git checkout`.",
						Line:    name.Token.Line,
						Column:  name.Token.Column,
					})
				}
			}
		}
	}

	return violations
}

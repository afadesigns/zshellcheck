package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1058",
		Title: "Avoid `sudo` with redirection",
		Description: "Redirecting output of `sudo` (e.g. `sudo cmd > /file`) fails if the current user " +
			"doesn't have permission. Use `| sudo tee /file` instead.",
		Check: checkZC1058,
	})
}

func checkZC1058(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	if name, ok := cmd.Name.(*ast.Identifier); !ok || name.Value != "sudo" {
		return nil
	}

	violations := []Violation{}

	for _, arg := range cmd.Arguments {
		argStr := arg.String()
		argStr = strings.Trim(argStr, "\"'")
		if argStr == ">" || argStr == ">>" {
			violations = append(violations, Violation{
				KataID:  "ZC1058",
				Message: "Redirecting `sudo` output happens as the current user. Use `| sudo tee file` to write with privileges.",
				Line:    arg.TokenLiteralNode().Line,
				Column:  arg.TokenLiteralNode().Column,
			})
		}
	}

	return violations
}

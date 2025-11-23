package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.RedirectionNode, Kata{
		ID:    "ZC1058",
		Title: "Avoid `sudo` with redirection",
		Description: "Redirecting output of `sudo` (e.g. `sudo cmd > /file`) fails if the current user " +
			"doesn't have permission. Use `| sudo tee /file` instead.",
		Check: checkZC1058,
	})
}

func checkZC1058(node ast.Node) []Violation {
	redir, ok := node.(*ast.Redirection)
	if !ok {
		return nil
	}

	// Check if Left side is a sudo command
	// Left is Expression.
	// Should be SimpleCommand.

	cmd, ok := redir.Left.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	if name, ok := cmd.Name.(*ast.Identifier); ok && name.Value == "sudo" {
		// Check operator: > or >> (output)
		if redir.Operator == ">" || redir.Operator == ">>" {
			return []Violation{{
				KataID:  "ZC1058",
				Message: "Redirecting `sudo` output happens as the current user. Use `| sudo tee file` to write with privileges.",
				Line:    redir.TokenLiteralNode().Line,
				Column:  redir.TokenLiteralNode().Column,
			}}
		}
	}

	return nil
}

package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:          "ZC1061",
		Title:       "Prefer `{start..end}` over `seq`",
		Description: "Using `seq` creates an external process. Zsh supports integer range expansion natively: `{1..10}`.",
		Check:       checkZC1061,
	})
}

func checkZC1061(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	// Check if command is seq
	if name, ok := cmd.Name.(*ast.Identifier); ok && name.Value == "seq" {
		return []Violation{{
			KataID:  "ZC1061",
			Message: "Prefer `{start..end}` range expansion over `seq`. It is built-in and faster.",
			Line:    name.Token.Line,
			Column:  name.Token.Column,
		}}
	}

	return nil
}

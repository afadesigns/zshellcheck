package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1013",
		Title: "Use `((...))` for arithmetic operations instead of `let`",
		Description: "The `let` command is a shell builtin, but the `((...))` syntax is more portable " +
			"and generally preferred for arithmetic operations in Zsh.",
		Check: checkZC1013,
	})
}

func checkZC1013(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	if cmd.Name.String() == "let" {
		return []Violation{{
			KataID:  "ZC1013",
			Message: "Use `((...))` for arithmetic operations instead of `let`.",
			Line:    cmd.TokenLiteralNode().Line,
			Column:  cmd.TokenLiteralNode().Column,
		}}
	}

	return nil
}

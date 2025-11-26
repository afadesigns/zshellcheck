package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1096",
		Title: "Warn on `bc` for simple arithmetic",
		Description: "Zsh has built-in support for floating point arithmetic using `(( ... ))` or `$(( ... ))`. " +
			"Using `bc` is often unnecessary and slower.",
		Check: checkZC1096,
	})
}

func checkZC1096(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	if cmd.Name != nil && cmd.Name.String() == "bc" {
		return []Violation{{
			KataID:  "ZC1096",
			Message: "Zsh supports floating point arithmetic natively. You often don't need `bc`.",
			Line:    cmd.Token.Line,
			Column:  cmd.Token.Column,
		}}
	}

	return nil
}

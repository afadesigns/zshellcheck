package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:          "ZC1063",
		Title:       "Prefer `grep -F` over `fgrep`",
		Description: "`fgrep` is deprecated. Use `grep -F` instead.",
		Check:       checkZC1063,
	})
}

func checkZC1063(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	if name, ok := cmd.Name.(*ast.Identifier); ok && name.Value == "fgrep" {
		return []Violation{{
			KataID:  "ZC1063",
			Message: "`fgrep` is deprecated. Use `grep -F` instead.",
			Line:    name.Token.Line,
			Column:  name.Token.Column,
		}}
	}

	return nil
}

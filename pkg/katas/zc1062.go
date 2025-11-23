package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:          "ZC1062",
		Title:       "Prefer `grep -E` over `egrep`",
		Description: "`egrep` is deprecated. Use `grep -E` instead.",
		Check:       checkZC1062,
	})
}

func checkZC1062(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	if name, ok := cmd.Name.(*ast.Identifier); ok && name.Value == "egrep" {
		return []Violation{{
			KataID:  "ZC1062",
			Message: "`egrep` is deprecated. Use `grep -E` instead.",
			Line:    name.Token.Line,
			Column:  name.Token.Column,
		}}
	}

	return nil
}

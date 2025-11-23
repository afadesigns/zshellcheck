package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:          "ZC1064",
		Title:       "Prefer `command -v` over `type`",
		Description: "`type` output format varies and is not POSIX standard for checking existence. " +
			"`command -v` is quieter and standard.",
		Check: checkZC1064,
	})
}

func checkZC1064(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	if name, ok := cmd.Name.(*ast.Identifier); ok && name.Value == "type" {
		return []Violation{{
			KataID:  "ZC1064",
			Message: "Prefer `command -v` over `type`. `type` output is not stable/standard for checking command existence.",
			Line:    name.Token.Line,
			Column:  name.Token.Column,
		}}
	}

	return nil
}

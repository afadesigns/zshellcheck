package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:          "ZC1049",
		Title:       "Prefer functions over aliases",
		Description: "Aliases are expanded at parse time and can be confusing in scripts. " +
			"Use functions for more predictable behavior.",
		Check: checkZC1049,
	})
}

func checkZC1049(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	// Check if command is alias
	if name, ok := cmd.Name.(*ast.Identifier); ok && name.Value == "alias" {
		return []Violation{{
			KataID:  "ZC1049",
			Message: "Prefer functions over aliases. " +
				"Aliases are expanded at parse time and can behave unexpectedly in scripts.",
			Line:    name.Token.Line,
			Column:  name.Token.Column,
		}}
	}

	return nil
}

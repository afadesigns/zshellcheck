package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1293",
		Title:    "Use `[[ ]]` instead of `test` command in Zsh",
		Severity: SeverityStyle,
		Description: "Zsh `[[ ]]` provides a more powerful conditional expression syntax than " +
			"the `test` command. It supports pattern matching, regex, and does not require " +
			"quoting of variable expansions to prevent word splitting.",
		Check: checkZC1293,
	})
}

func checkZC1293(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "test" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1293",
		Message: "Use `[[ ]]` instead of the `test` command in Zsh. `[[ ]]` is more powerful and does not require variable quoting.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityStyle,
	}}
}

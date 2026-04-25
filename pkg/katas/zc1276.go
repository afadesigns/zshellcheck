package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1276",
		Title:    "Use Zsh `{start..end}` instead of `seq`",
		Severity: SeverityStyle,
		Description: "Zsh natively supports `{start..end}` brace expansion for generating number " +
			"sequences, avoiding the overhead of forking the external `seq` command.",
		Check: checkZC1276,
		Fix:   fixZC1061,
	})
}

func checkZC1276(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "seq" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1276",
		Message: "Use Zsh brace expansion `{start..end}` instead of `seq`. Brace expansion is built-in and avoids forking.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityStyle,
	}}
}

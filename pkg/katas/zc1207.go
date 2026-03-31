package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1207",
		Title:    "Avoid `passwd` in scripts — use `chpasswd`",
		Severity: SeverityWarning,
		Description: "`passwd` prompts interactively for password input. " +
			"Use `chpasswd` or `usermod --password` for non-interactive password changes.",
		Check: checkZC1207,
	})
}

func checkZC1207(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "passwd" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1207",
		Message: "Avoid `passwd` in scripts — it prompts interactively. " +
			"Use `chpasswd` or `usermod --password` for non-interactive password changes.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

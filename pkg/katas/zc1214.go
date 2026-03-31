package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1214",
		Title:    "Avoid `su` in scripts — use `sudo -u` for user switching",
		Severity: SeverityWarning,
		Description: "`su` prompts for a password interactively which hangs scripts. " +
			"Use `sudo -u user cmd` for non-interactive privilege switching.",
		Check: checkZC1214,
	})
}

func checkZC1214(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "su" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1214",
		Message: "Avoid `su` in scripts — it prompts for a password interactively. " +
			"Use `sudo -u user cmd` for non-interactive privilege switching.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1195",
		Title:    "Avoid overly permissive `umask` values",
		Severity: SeverityWarning,
		Description: "`umask 000` or `umask 0000` creates world-writable files by default. " +
			"Use `umask 022` or more restrictive values for security.",
		Check: checkZC1195,
	})
}

func checkZC1195(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "umask" {
		return nil
	}

	if len(cmd.Arguments) != 1 {
		return nil
	}

	val := cmd.Arguments[0].String()
	if val == "000" || val == "0000" || val == "0" {
		return []Violation{{
			KataID: "ZC1195",
			Message: "Avoid `umask 000` — it creates world-writable files. " +
				"Use `umask 022` or `umask 077` for secure default permissions.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}

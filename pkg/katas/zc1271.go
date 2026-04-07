package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1271",
		Title:    "Use `command -v` instead of `which` for command existence checks",
		Severity: SeverityStyle,
		Description: "`which` is not POSIX-standard and behaves inconsistently across systems. " +
			"Use `command -v` which is portable and built into Zsh.",
		Check: checkZC1271,
	})
}

func checkZC1271(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "which" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1271",
		Message: "Use `command -v` instead of `which`. `command -v` is POSIX-compliant and built into Zsh.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityStyle,
	}}
}

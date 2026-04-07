package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1278",
		Title:    "Use `$(( ))` instead of `expr` for arithmetic",
		Severity: SeverityStyle,
		Description: "`expr` is an external command for arithmetic. Zsh has native arithmetic " +
			"expansion `$(( ))` which is faster and more readable.",
		Check: checkZC1278,
	})
}

func checkZC1278(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "expr" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1278",
		Message: "Use Zsh arithmetic expansion `$(( ))` instead of `expr`. It is built-in and avoids forking.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityStyle,
	}}
}

package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1160",
		Title:    "Prefer `curl` over `wget` for portability",
		Severity: SeverityStyle,
		Description: "`wget` is not installed by default on macOS. " +
			"`curl` is available on virtually all Unix systems and is more portable.",
		Check: checkZC1160,
	})
}

func checkZC1160(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "wget" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1160",
		Message: "Prefer `curl` over `wget` for portability. " +
			"`curl` is pre-installed on macOS and most Linux distributions.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

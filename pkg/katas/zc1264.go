package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1264",
		Title:    "Use `dnf` instead of `yum` on modern Fedora/RHEL",
		Severity: SeverityStyle,
		Description: "`yum` is deprecated on Fedora 22+ and RHEL 8+. " +
			"`dnf` is the modern replacement with better dependency resolution.",
		Check: checkZC1264,
	})
}

func checkZC1264(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "yum" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1264",
		Message: "Use `dnf` instead of `yum`. `yum` is deprecated on modern " +
			"Fedora and RHEL; `dnf` has better dependency resolution.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

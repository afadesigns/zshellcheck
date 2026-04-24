package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1263",
		Title:    "Use `apt-get` instead of `apt` in scripts",
		Severity: SeverityStyle,
		Description: "`apt` is designed for interactive use and its output format may change. " +
			"`apt-get` has a stable interface suitable for scripts and CI.",
		Check: checkZC1263,
		Fix:   fixZC1263,
	})
}

// fixZC1263 rewrites `apt` to `apt-get` at the command name position.
// Arguments stay intact — the two tools accept the same shape for
// install / upgrade / update / remove.
func fixZC1263(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "apt" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("apt"),
		Replace: "apt-get",
	}}
}

func checkZC1263(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "apt" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1263",
		Message: "Use `apt-get` instead of `apt` in scripts. " +
			"`apt` is for interactive use; `apt-get` has a stable scripting interface.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

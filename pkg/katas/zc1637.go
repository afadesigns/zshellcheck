package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1637",
		Title:    "Style: prefer Zsh `typeset -r NAME=value` over POSIX `readonly NAME=value`",
		Severity: SeverityStyle,
		Description: "Both `readonly NAME` and `typeset -r NAME` create a read-only parameter. " +
			"In Zsh the idiomatic form is `typeset -r` — it composes with other typeset flags " +
			"(`-ir` for readonly integer, `-xr` for readonly export, `-gr` to pin a readonly " +
			"global from inside a function). `readonly` works but reads as a Bash / POSIX-ism " +
			"in a Zsh codebase.",
		Check: checkZC1637,
		Fix:   fixZC1637,
	})
}

// fixZC1637 rewrites the `readonly` command name to `typeset -r`.
// Single-edit replacement at the violation column. Detector gates on
// the bare command name match, so the rewrite is idempotent on a
// re-run (the new line starts with `typeset` not `readonly`).
func fixZC1637(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "readonly" {
		return nil
	}
	off := LineColToByteOffset(source, v.Line, v.Column)
	if off < 0 || off+len("readonly") > len(source) {
		return nil
	}
	if string(source[off:off+len("readonly")]) != "readonly" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("readonly"),
		Replace: "typeset -r",
	}}
}

func checkZC1637(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "readonly" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}

	return []Violation{{
		KataID: "ZC1637",
		Message: "`readonly` works but `typeset -r NAME=value` is the Zsh-native form and " +
			"composes with other typeset flags.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

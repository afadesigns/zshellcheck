// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1335",
		Title:    "Use Zsh array reversal instead of `tac` for in-memory data",
		Severity: SeverityStyle,
		Description: "`tac` reverses lines from a file or stdin. For in-memory array data, " +
			"Zsh provides `${(Oa)array}` to reverse array element order without " +
			"spawning an external process.",
		Check: checkZC1335,
	})
}

func checkZC1335(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "tac" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1335",
		Message: "Consider Zsh `${(Oa)array}` for reversing array data instead of piping to `tac`.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityStyle,
	}}
}

// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1294",
		Title:    "Use `bindkey` instead of `bind` for key bindings in Zsh",
		Severity: SeverityWarning,
		Description: "`bind` is a Bash builtin for key bindings. Zsh uses `bindkey` for " +
			"ZLE (Zsh Line Editor) key bindings. Using `bind` in a Zsh script will " +
			"fail unless Bash compatibility is loaded.",
		Check: checkZC1294,
	})
}

func checkZC1294(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "bind" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1294",
		Message: "Use `bindkey` instead of `bind` in Zsh. `bind` is a Bash builtin; Zsh uses `bindkey` for ZLE key bindings.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityWarning,
	}}
}

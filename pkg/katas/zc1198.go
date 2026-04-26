// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1198",
		Title:    "Avoid interactive editors in scripts",
		Severity: SeverityWarning,
		Description: "`nano`, `vi`, and `vim` require interactive terminals and will hang " +
			"in non-interactive scripts. Use `sed -i` or `ed` for scripted editing.",
		Check: checkZC1198,
	})
}

func checkZC1198(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	name := ident.Value
	if name != "nano" && name != "vi" && name != "vim" && name != "emacs" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1198",
		Message: "Avoid `" + name + "` in scripts — interactive editors hang without a terminal. " +
			"Use `sed -i` or `ed` for scripted file editing.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

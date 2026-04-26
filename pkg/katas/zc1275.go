// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1275",
		Title:    "Use Zsh `${var:h}` instead of `dirname`",
		Severity: SeverityStyle,
		Description: "Zsh provides the `:h` (head) modifier for parameter expansion which extracts " +
			"the directory component, avoiding the overhead of forking `dirname`.",
		Check: checkZC1275,
	})
}

func checkZC1275(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "dirname" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1275",
		Message: "Use Zsh parameter expansion `${var:h}` instead of `dirname`. The `:h` modifier extracts the directory without forking a process.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityStyle,
	}}
}

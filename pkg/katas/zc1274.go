// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1274",
		Title:    "Use Zsh `${var:t}` instead of `basename`",
		Severity: SeverityStyle,
		Description: "Zsh provides the `:t` (tail) modifier for parameter expansion which extracts " +
			"the filename component, avoiding the overhead of forking `basename`.",
		Check: checkZC1274,
	})
}

func checkZC1274(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "basename" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1274",
		Message: "Use Zsh parameter expansion `${var:t}` instead of `basename`. The `:t` modifier extracts the filename without forking a process.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityStyle,
	}}
}

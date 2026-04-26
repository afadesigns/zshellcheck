// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1256",
		Title:    "Clean up `mkfifo` pipes with a trap on EXIT",
		Severity: SeverityInfo,
		Description: "`mkfifo` creates named pipes that persist on the filesystem. " +
			"Set up a `trap` to remove them on EXIT to prevent leftover files.",
		Check: checkZC1256,
	})
}

func checkZC1256(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "mkfifo" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1256",
		Message: "Set up `trap 'rm -f pipe' EXIT` after `mkfifo`. " +
			"Named pipes persist on the filesystem and need explicit cleanup.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}

// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1370",
		Title:    "Prefer Zsh `repeat N { ... }` over `yes str | head -n N` for finite output",
		Severity: SeverityStyle,
		Description: "`yes` plus `head` is a common idiom for producing N copies of a line. " +
			"Zsh's `repeat N { print str }` does the same loop in-shell without spawning yes " +
			"or the pipe, and without the SIGPIPE handshake.",
		Check: checkZC1370,
	})
}

func checkZC1370(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "yes" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1370",
		Message: "Prefer Zsh `repeat N { print str }` over `yes str | head -n N` for producing " +
			"N copies of a line. No external `yes` process, no pipe.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

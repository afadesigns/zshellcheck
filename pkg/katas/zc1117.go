// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1117",
		Title: "Use `&!` or `disown` instead of `nohup`",
		Description: "Zsh provides `&!` (shorthand for `& disown`) to run a command in the background " +
			"immune to hangups. Avoid spawning `nohup` as an external process.",
		Severity: SeverityStyle,
		Check:    checkZC1117,
	})
}

func checkZC1117(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "nohup" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1117",
		Message: "Use `cmd &!` or `cmd & disown` instead of `nohup cmd &`. " +
			"Zsh `&!` is a built-in shorthand that avoids spawning nohup.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

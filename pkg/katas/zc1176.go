// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1176",
		Title:    "Use `zparseopts` instead of `getopt`/`getopts`",
		Severity: SeverityStyle,
		Description: "Zsh provides `zparseopts` for powerful option parsing with long options, " +
			"arrays, and defaults. Avoid `getopt`/`getopts` which are less capable in Zsh.",
		Check: checkZC1176,
	})
}

func checkZC1176(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "getopt" && ident.Value != "getopts" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1176",
		Message: "Use Zsh `zparseopts` instead of `" + ident.Value + "`. " +
			"`zparseopts` supports long options, arrays, and is the native Zsh approach.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

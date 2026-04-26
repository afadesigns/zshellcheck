// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1123",
		Title: "Use `$OSTYPE` instead of `uname`",
		Description: "Zsh provides `$OSTYPE` (e.g., `linux-gnu`, `darwin`) as a built-in variable. " +
			"Avoid spawning `uname` for simple OS detection.",
		Severity: SeverityStyle,
		Check:    checkZC1123,
	})
}

func checkZC1123(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "uname" {
		return nil
	}

	// Only flag simple uname, uname -s, uname -o (OS type detection)
	// Skip uname -r, -m, -a, -n, -p which provide different info
	if len(cmd.Arguments) == 0 {
		return []Violation{{
			KataID: "ZC1123",
			Message: "Use `$OSTYPE` instead of `uname` for OS detection. " +
				"Zsh maintains `$OSTYPE` as a built-in variable.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-s" || val == "-o" {
			return []Violation{{
				KataID: "ZC1123",
				Message: "Use `$OSTYPE` instead of `uname -s` for OS detection. " +
					"Zsh maintains `$OSTYPE` as a built-in variable.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

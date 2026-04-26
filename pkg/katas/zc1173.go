// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1173",
		Title:    "Avoid `column` command — use Zsh `print -C` for columnar output",
		Severity: SeverityStyle,
		Description: "Zsh `print -C N` formats output into N columns natively. " +
			"Avoid spawning `column` as an external process for simple tabulation.",
		Check: checkZC1173,
	})
}

func checkZC1173(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "column" {
		return nil
	}

	// Only flag simple column usage (column -t is the most common)
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-t" {
			return []Violation{{
				KataID: "ZC1173",
				Message: "Use Zsh `print -C N` for columnar output instead of `column -t`. " +
					"The `print` builtin formats columns without spawning an external process.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

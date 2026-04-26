// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1179",
		Title:    "Use Zsh `strftime` instead of `date` for formatting",
		Severity: SeverityStyle,
		Description: "Zsh provides `strftime` via `zmodload zsh/datetime` for date formatting. " +
			"Avoid spawning `date` for simple timestamp formatting.",
		Check: checkZC1179,
	})
}

func checkZC1179(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "date" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := strings.Trim(arg.String(), "'\"")
		if len(val) > 1 && val[0] == '+' && val != "+%s" && val != "+%s%N" {
			return []Violation{{
				KataID: "ZC1179",
				Message: "Use `strftime` (via `zmodload zsh/datetime`) instead of `date +" + val[1:] + "`. " +
					"Zsh date formatting avoids spawning an external process.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

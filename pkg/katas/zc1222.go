// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1222",
		Title:    "Avoid `lsof -i` for port checks — use `ss -tlnp`",
		Severity: SeverityStyle,
		Description: "`lsof -i` is slow and requires elevated permissions on some systems. " +
			"`ss -tlnp` is faster and part of the standard iproute2 toolkit.",
		Check: checkZC1222,
	})
}

func checkZC1222(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "lsof" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-i" {
			return []Violation{{
				KataID: "ZC1222",
				Message: "Use `ss -tlnp` instead of `lsof -i` for port checks. " +
					"`ss` is faster and doesn't require elevated permissions.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

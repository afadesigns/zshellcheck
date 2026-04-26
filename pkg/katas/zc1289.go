// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1289",
		Title:    "Use Zsh `${(u)array}` for unique elements instead of `sort -u`",
		Severity: SeverityStyle,
		Description: "Zsh provides the `(u)` parameter expansion flag to remove duplicate " +
			"elements from an array. This preserves original order and avoids spawning " +
			"an external `sort -u` process.",
		Check: checkZC1289,
	})
}

func checkZC1289(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "sort" {
		return nil
	}

	hasUnique := false
	hasOtherFlags := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-u" {
			hasUnique = true
		} else if len(val) > 1 && val[0] == '-' {
			hasOtherFlags = true
		}
	}

	if hasUnique && !hasOtherFlags {
		return []Violation{{
			KataID:  "ZC1289",
			Message: "Use Zsh `${(u)array}` for unique elements instead of `sort -u`. The `(u)` flag preserves order.",
			Line:    cmd.Token.Line,
			Column:  cmd.Token.Column,
			Level:   SeverityStyle,
		}}
	}

	return nil
}

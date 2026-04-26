// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1291",
		Title:    "Use Zsh `${(O)array}` for reverse sorting instead of `sort -r`",
		Severity: SeverityStyle,
		Description: "Zsh provides the `(O)` parameter expansion flag to sort array elements " +
			"in descending (reverse) order. This avoids spawning an external `sort -r` " +
			"process for simple reverse sorting of array data.",
		Check: checkZC1291,
	})
}

func checkZC1291(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "sort" {
		return nil
	}

	hasReverse := false
	hasOtherFlags := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-r" {
			hasReverse = true
		} else if len(val) > 1 && val[0] == '-' {
			hasOtherFlags = true
		}
	}

	if hasReverse && !hasOtherFlags {
		return []Violation{{
			KataID:  "ZC1291",
			Message: "Use Zsh `${(O)array}` for reverse sorting instead of `sort -r`. The `(O)` flag sorts descending in-shell.",
			Line:    cmd.Token.Line,
			Column:  cmd.Token.Column,
			Level:   SeverityStyle,
		}}
	}

	return nil
}

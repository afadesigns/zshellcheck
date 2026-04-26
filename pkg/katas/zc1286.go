// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1286",
		Title:    "Use Zsh `${array:#pattern}` instead of `grep -v` for filtering",
		Severity: SeverityStyle,
		Description: "Zsh provides `${array:#pattern}` to remove matching elements from an array " +
			"and `${(M)array:#pattern}` to keep only matching elements. This avoids " +
			"spawning an external `grep` process for simple filtering tasks.",
		Check: checkZC1286,
	})
}

func checkZC1286(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "grep" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-v" {
			return []Violation{{
				KataID:  "ZC1286",
				Message: "Use Zsh `${array:#pattern}` for filtering instead of `grep -v`. Parameter expansion avoids a subprocess.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
				Level:   SeverityStyle,
			}}
		}
	}

	return nil
}

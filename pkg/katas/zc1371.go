// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1371",
		Title:    "Use Zsh array `:t` modifier instead of `basename -a` for bulk path stripping",
		Severity: SeverityStyle,
		Description: "`basename -a a b c` returns the file name component of each path. Zsh's " +
			"`${array:t}` parameter modifier applies the same tail-component extraction to every " +
			"element of an array at once — no external process.",
		Check: checkZC1371,
	})
}

func checkZC1371(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "basename" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-a" || v == "--multiple" {
			return []Violation{{
				KataID: "ZC1371",
				Message: "Use Zsh `${paths:t}` on an array for bulk basename extraction instead of " +
					"`basename -a`. The `:t` modifier applies to every array element.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1338",
		Title:    "Avoid `seq -s` — use Zsh `${(j:sep:)${(s::)...}}` for joining",
		Severity: SeverityStyle,
		Description: "`seq -s` generates a sequence with a custom separator. Zsh provides " +
			"native brace expansion with `{start..end}` and `${(j:sep:)array}` " +
			"for joining, avoiding an external process.",
		Check: checkZC1338,
	})
}

func checkZC1338(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "seq" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-s" {
			return []Violation{{
				KataID:  "ZC1338",
				Message: "Avoid `seq -s` in Zsh — use `${(j:sep:)array}` with brace expansion for joined sequences.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
				Level:   SeverityStyle,
			}}
		}
	}

	return nil
}

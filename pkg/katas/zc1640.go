// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1640",
		Title:    "Style: `${!var}` Bash indirect expansion — prefer Zsh `${(P)var}`",
		Severity: SeverityStyle,
		Description: "`${!var}` is Bash indirect expansion — it reads the value of the " +
			"parameter whose name is stored in `$var`. Zsh has the native flag form " +
			"`${(P)var}` which does the same and composes with other parameter-expansion " +
			"flags (`${(Pf)var}` to split the indirect value on newlines, for example). " +
			"`${!prefix*}` / `${!array[@]}` have Zsh equivalents via the `$parameters` hash " +
			"or `(k)` subscript flags. Prefer the native Zsh form in a Zsh codebase.",
		Check: checkZC1640,
	})
}

func checkZC1640(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.Contains(v, "${!") {
			return []Violation{{
				KataID: "ZC1640",
				Message: "`${!var}` Bash indirect — prefer Zsh `${(P)var}` for the same " +
					"semantics with flag composability.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}
	return nil
}

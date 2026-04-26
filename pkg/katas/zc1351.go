// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1351",
		Title:    "Use `[[ $str =~ pattern ]]` instead of `expr match` / `expr :` for regex",
		Severity: SeverityStyle,
		Description: "Zsh's `[[ $str =~ pattern ]]` evaluates regex natively and populates `$match` / " +
			"`$MATCH` / `$mbegin` / `$mend` arrays. Avoid shelling out to `expr match` or the " +
			"`expr STRING : REGEX` form.",
		Check: checkZC1351,
	})
}

func checkZC1351(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "expr" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "match" || v == "index" {
			return []Violation{{
				KataID: "ZC1351",
				Message: "Use `[[ $str =~ pattern ]]` with `$match` / `$MATCH` arrays instead of " +
					"`expr match`/`expr index`. Regex evaluation stays in the shell.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

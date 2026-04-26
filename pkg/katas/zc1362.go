// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1362",
		Title:    "Use `[[ -o option ]]` instead of `test -o option` for Zsh option checks",
		Severity: SeverityInfo,
		Description: "In Zsh, `[[ -o name ]]` tests whether a shell option is set. The `test` / `[` " +
			"builtin interprets `-o` as a logical OR, not an option-query — so `test -o foo` is " +
			"a syntax error or wrong behavior. Use the `[[ ... ]]` form for option tests.",
		Check: checkZC1362,
	})
}

func checkZC1362(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "test" && ident.Value != "[" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-o" {
			return []Violation{{
				KataID: "ZC1362",
				Message: "Use `[[ -o option ]]` for option checks in Zsh — `test -o` means logical OR, " +
					"not option-query, producing wrong results.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityInfo,
			}}
		}
	}

	return nil
}

// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1354",
		Title:    "Use `whence -w` instead of Bash-specific `type -t` for command classification",
		Severity: SeverityStyle,
		Description: "`type -t` returns the category (alias, keyword, function, builtin, file) " +
			"of a command in Bash. Zsh's `whence -w` produces `name: category` output with " +
			"the same information and without shelling out for the sub-field extraction.",
		Check: checkZC1354,
	})
}

func checkZC1354(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "type" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-t" || v == "-a" || v == "-P" {
			return []Violation{{
				KataID: "ZC1354",
				Message: "Use Zsh `whence -w` (category), `whence -a` (all), or `whence -p` (path) " +
					"instead of Bash-specific `type -t`/`-a`/`-P`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

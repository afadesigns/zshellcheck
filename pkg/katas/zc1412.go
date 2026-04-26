// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1412",
		Title:    "Avoid `$COMPREPLY` — Bash completion output, use Zsh `compadd`",
		Severity: SeverityError,
		Description: "Bash completion functions populate the `$COMPREPLY` array to declare " +
			"candidates. Zsh's compsys uses the `compadd` builtin: `compadd -- foo bar baz`. " +
			"Setting `$COMPREPLY` in a Zsh completion does nothing.",
		Check: checkZC1412,
	})
}

func checkZC1412(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "echo" && ident.Value != "print" && ident.Value != "printf" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.Contains(v, "COMPREPLY") {
			return []Violation{{
				KataID: "ZC1412",
				Message: "`$COMPREPLY` is a Bash-only completion output array. In Zsh compsys " +
					"use `compadd -- candidate1 candidate2`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}

	return nil
}

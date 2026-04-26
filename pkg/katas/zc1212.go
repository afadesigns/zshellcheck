// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1212",
		Title:    "Avoid `git add .` — use explicit paths or `git add -p`",
		Severity: SeverityInfo,
		Description: "`git add .` stages everything including unintended files. " +
			"Use explicit file paths or `git add -p` for selective staging.",
		Check: checkZC1212,
	})
}

func checkZC1212(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "git" {
		return nil
	}

	if len(cmd.Arguments) < 2 {
		return nil
	}

	if cmd.Arguments[0].String() != "add" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		val := arg.String()
		if val == "." || val == "-A" {
			return []Violation{{
				KataID: "ZC1212",
				Message: "Avoid `git add .` or `git add -A` — they stage everything including " +
					"unintended files. Use explicit paths or `git add -p` for selective staging.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityInfo,
			}}
		}
	}

	return nil
}

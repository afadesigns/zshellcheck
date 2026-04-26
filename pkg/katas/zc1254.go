// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1254",
		Title:    "Avoid `git commit --amend` in shared branches",
		Severity: SeverityWarning,
		Description: "`git commit --amend` rewrites the last commit which causes problems " +
			"if already pushed. Use `git commit --fixup` or a new commit instead.",
		Check: checkZC1254,
	})
}

func checkZC1254(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "git" {
		return nil
	}

	if len(cmd.Arguments) < 1 || cmd.Arguments[0].String() != "commit" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		val := arg.String()
		if val == "--amend" {
			return []Violation{{
				KataID: "ZC1254",
				Message: "Avoid `git commit --amend` on shared branches — it rewrites history. " +
					"Use `git commit --fixup` or create a new commit instead.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}

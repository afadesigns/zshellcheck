// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1184",
		Title:    "Avoid `diff -u` for patch generation — use `git diff` when in a repo",
		Severity: SeverityStyle,
		Description: "When working within a git repository, `git diff` provides better context, " +
			"color output, and integration. Use `diff -u` only for non-repo file comparisons.",
		Check: checkZC1184,
	})
}

func checkZC1184(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "diff" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-u" || val == "--unified" {
			return []Violation{{
				KataID: "ZC1184",
				Message: "Consider `git diff` instead of `diff -u` when working in a repository. " +
					"`git diff` provides better context and integration.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

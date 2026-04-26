// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1719",
		Title:    "Warn on `git filter-branch` — deprecated since Git 2.24, use `git filter-repo`",
		Severity: SeverityWarning,
		Description: "`git filter-branch` is deprecated as of Git 2.24; its manpage opens with " +
			"\"WARNING: this command is deprecated\" and points users at `git filter-repo`. " +
			"`filter-branch` is single-process slow, mishandles common cases (tag rewrites, " +
			"refs/notes/*, signed commits), and leaves orphaned objects behind. The modern " +
			"replacement is `git filter-repo` (separate package; `apt/brew install git-" +
			"filter-repo`) — much faster, safer defaults, and what GitHub / GitLab guidance " +
			"recommends for secret-removal rewrites.",
		Check: checkZC1719,
	})
}

func checkZC1719(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "git" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	if cmd.Arguments[0].String() != "filter-branch" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1719",
		Message: "`git filter-branch` is deprecated (Git 2.24+) and its manpage redirects to " +
			"`git filter-repo`. Use that instead — faster, safer defaults, no orphaned " +
			"objects.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

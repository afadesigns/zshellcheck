// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1714",
		Title:    "Error on `gh repo delete --yes` / `gh release delete --yes` — bypassed confirmation",
		Severity: SeverityError,
		Description: "`gh repo delete OWNER/REPO --yes` (and `gh release delete TAG --yes`) " +
			"skip the interactive confirmation that protects against typos and broken " +
			"variable expansion. A repository deletion is final — issues, PRs, releases, " +
			"GitHub Actions history, and (for free accounts) any forks against it all " +
			"disappear with no soft-delete window. From a script, run without `--yes` so a " +
			"human reviews the target, or wrap deletion in a manually-triggered workflow " +
			"with explicit input prompts.",
		Check: checkZC1714,
	})
}

func checkZC1714(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "gh" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	scope := cmd.Arguments[0].String()
	if scope != "repo" && scope != "release" {
		return nil
	}
	if cmd.Arguments[1].String() != "delete" {
		return nil
	}

	for _, arg := range cmd.Arguments[2:] {
		if arg.String() == "--yes" {
			return []Violation{{
				KataID: "ZC1714",
				Message: "`gh " + scope + " delete --yes` bypasses GitHub's confirmation — " +
					"a typo or stale variable destroys the target with no soft-delete. " +
					"Drop `--yes` so a human confirms.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

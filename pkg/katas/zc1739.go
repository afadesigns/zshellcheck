// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1739",
		Title:    "Warn on `git submodule update --remote` — pulls upstream HEAD, breaks reproducibility",
		Severity: SeverityWarning,
		Description: "`git submodule update --remote` fetches each submodule's tracked branch HEAD " +
			"instead of the commit pinned in the parent repo's index. Builds become " +
			"non-reproducible — every CI run pulls a different commit — and any compromised " +
			"upstream commit lands directly in the build. Use `git submodule update --init " +
			"--recursive` (defaults to the pinned commit) and bump submodule pins through " +
			"reviewed PRs.",
		Check: checkZC1739,
	})
}

func checkZC1739(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "git" {
		return nil
	}
	if len(cmd.Arguments) < 3 {
		return nil
	}
	if cmd.Arguments[0].String() != "submodule" || cmd.Arguments[1].String() != "update" {
		return nil
	}

	for _, arg := range cmd.Arguments[2:] {
		if arg.String() == "--remote" {
			return []Violation{{
				KataID: "ZC1739",
				Message: "`git submodule update --remote` ignores the pinned commits in the " +
					"parent repo and pulls each submodule's branch HEAD — non-" +
					"reproducible builds, supply-chain risk. Use `--init --recursive` " +
					"and bump pins via reviewed PRs.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

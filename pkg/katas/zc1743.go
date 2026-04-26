// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1743",
		Title:    "Warn on `npm audit fix --force` — accepts major-version dependency bumps silently",
		Severity: SeverityWarning,
		Description: "`npm audit fix --force` (and `pnpm audit --fix --force`) resolves advisories " +
			"by upgrading dependencies past semver-major boundaries when no backward-" +
			"compatible patch exists. The flag accepts every upgrade without surfacing the " +
			"breaking changes — a build can silently move to a new major of a transitive " +
			"dependency that removes APIs your code calls. Drop `--force` and triage each " +
			"advisory individually; `npm audit fix` handles compatible patches, and the " +
			"remaining advisory targets need a pin or a vendored fork.",
		Check: checkZC1743,
	})
}

func checkZC1743(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var matches bool
	switch ident.Value {
	case "npm":
		if len(cmd.Arguments) >= 2 && cmd.Arguments[0].String() == "audit" && cmd.Arguments[1].String() == "fix" {
			matches = true
		}
	case "pnpm":
		if len(cmd.Arguments) >= 2 && cmd.Arguments[0].String() == "audit" {
			hasFix := false
			for _, a := range cmd.Arguments[1:] {
				if a.String() == "--fix" {
					hasFix = true
					break
				}
			}
			if hasFix {
				matches = true
			}
		}
	default:
		return nil
	}
	if !matches {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "--force" || arg.String() == "-f" {
			return []Violation{{
				KataID: "ZC1743",
				Message: "`" + ident.Value + " audit ... --force` accepts every major-" +
					"version bump an advisory triggers — silent breaking changes. Drop " +
					"`--force` and triage advisories one by one.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

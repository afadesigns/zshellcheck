// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1504",
		Title:    "Warn on `git push --mirror` — overwrites every remote ref",
		Severity: SeverityWarning,
		Description: "`git push --mirror` pushes every ref under `refs/` and deletes any remote " +
			"ref that is not present locally. Running it against a shared origin instantly " +
			"wipes everyone else's branches and tags. Legitimate uses are mirror-to-mirror " +
			"replication where the source is the authoritative tree; for everyday pushes use " +
			"an explicit refspec or `git push --all`.",
		Check: checkZC1504,
	})
}

func checkZC1504(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "git" {
		return nil
	}

	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "push" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "--mirror" {
			return []Violation{{
				KataID: "ZC1504",
				Message: "`git push --mirror` overwrites every remote ref and deletes ones " +
					"missing locally. Use an explicit refspec or `--all` for everyday pushes.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

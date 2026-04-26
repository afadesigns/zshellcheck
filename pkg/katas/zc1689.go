// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1689",
		Title:    "Error on `borg delete --force` — forced deletion of backup archives or repository",
		Severity: SeverityError,
		Description: "`borg delete --force REPO[::ARCHIVE]` bypasses the confirmation prompt " +
			"and removes the archive (or the whole repository, if ARCHIVE is omitted) in " +
			"one go. Unlike `borg prune`, which keeps a retention ladder and logs what it " +
			"would drop, `--force` deletion leaves nothing to restore from if the target " +
			"was typed wrong. Keep scripts to `borg prune --keep-daily` / `--keep-within` " +
			"with an explicit retention policy and gate any outright `borg delete` behind " +
			"a human `--checkpoint-interval` review.",
		Check: checkZC1689,
	})
}

func checkZC1689(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "borg" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "delete" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		if arg.String() == "--force" {
			return []Violation{{
				KataID: "ZC1689",
				Message: "`borg delete --force` skips confirmation and can nuke the whole " +
					"repository on a typo — use `borg prune --keep-*` with a retention " +
					"policy, or gate outright deletion behind a manual review.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

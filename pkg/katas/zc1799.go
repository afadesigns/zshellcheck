// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1799",
		Title:    "Warn on `rclone sync SRC DST` without `--dry-run` — one-way mirror can wipe DST",
		Severity: SeverityWarning,
		Description: "`rclone sync` makes DST look exactly like SRC: anything in DST that isn't in " +
			"SRC is deleted, including object versions on providers that support them. If SRC " +
			"is accidentally empty (typo in path, unmounted drive, wrong credentials " +
			"pointing at an empty bucket), the command silently wipes every object under DST " +
			"without a confirmation prompt. Always preview the diff with `rclone sync " +
			"--dry-run SRC DST` first; when you commit to the sync, keep `--backup-dir`, " +
			"`--max-delete`, or `--min-age` guards so a bad SRC cannot cascade into " +
			"unbounded deletion.",
		Check: checkZC1799,
	})
}

func checkZC1799(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "rclone" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "sync" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "--dry-run" || v == "-n" || v == "--interactive" || v == "-i" {
			return nil
		}
	}
	return []Violation{{
		KataID: "ZC1799",
		Message: "`rclone sync` deletes anything in DST that's not in SRC — empty / " +
			"wrong SRC silently wipes DST. Preview with `rclone sync --dry-run`, and " +
			"pin guards like `--backup-dir`, `--max-delete`, or `--min-age`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

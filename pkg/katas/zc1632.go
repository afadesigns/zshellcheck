// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1632",
		Title:    "Warn on `shred` — unreliable on journaled / CoW filesystems (ext4, btrfs, zfs)",
		Severity: SeverityWarning,
		Description: "`shred` assumes in-place overwrites, which is how ext2 worked. On a " +
			"journaled ext4 the overwrite passes go through the journal and may not hit the " +
			"original data blocks. On CoW filesystems (btrfs, zfs, xfs with reflink) the " +
			"overwrite lands in fresh blocks and leaves the old content intact until garbage " +
			"collection decides otherwise. `shred`'s own man page warns about this. For modern " +
			"secure deletion, use full-disk encryption with key destruction, or retire the " +
			"device with `blkdiscard` on an SSD.",
		Check: checkZC1632,
	})
}

func checkZC1632(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "shred" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}

	return []Violation{{
		KataID: "ZC1632",
		Message: "`shred` may not overwrite original blocks on ext4/btrfs/zfs. For " +
			"guaranteed erasure, use full-disk encryption with key destruction, or " +
			"`blkdiscard` when retiring an SSD.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

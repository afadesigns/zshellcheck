// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1994",
		Title:    "Error on `lvreduce -f` / `lvreduce -y` — shrinks the LV without checking the filesystem above",
		Severity: SeverityError,
		Description: "`lvreduce -L SIZE $LV` cuts the block device below an existing filesystem. " +
			"The confirmation prompt exists precisely because ext4/xfs/btrfs do not " +
			"shrink themselves — LVM happily lops off the tail even though the " +
			"filesystem still believes those blocks are allocated. `-f` / `-y` / " +
			"`--force` / `--yes` skip the prompt, and the next mount returns " +
			"corruption or missing files. Shrink the filesystem first with " +
			"`resize2fs $LV $NEWSIZE` (or `xfs_growfs` equivalent — xfs cannot shrink, " +
			"so offline backup + recreate), verify `df` / `fsck`, then `lvreduce " +
			"--resizefs` (which performs both steps atomically) instead of bypassing " +
			"the prompt.",
		Check: checkZC1994,
	})
}

func checkZC1994(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "lvreduce" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "-f", "-y", "--force", "--yes":
			return []Violation{{
				KataID: "ZC1994",
				Message: "`lvreduce " + v + "` skips the shrink-confirmation prompt — " +
					"the filesystem above still believes the tail is allocated and " +
					"the next mount sees corruption. Shrink fs first, or use " +
					"`lvreduce --resizefs`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

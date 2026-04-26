// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1749",
		Title:    "Error on `virsh undefine DOMAIN --remove-all-storage` — wipes VM disk images",
		Severity: SeverityError,
		Description: "`virsh undefine DOMAIN --remove-all-storage` (also `--wipe-storage` and the " +
			"newer `--storage <vol,vol>`) removes the VM's configuration AND deletes every " +
			"disk image the domain references. There is no soft-delete and no recycle bin — " +
			"a misresolved DOMAIN or a shared storage pool turns one typo into data loss " +
			"across VMs that happened to share a snapshot chain. Split the operation: back " +
			"up the qcow2 images (`virsh vol-clone` or `qemu-img convert`), then `virsh " +
			"undefine` without the storage flags, then delete volumes deliberately with " +
			"`virsh vol-delete` after a review.",
		Check: checkZC1749,
	})
}

func checkZC1749(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "virsh" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "undefine" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "--remove-all-storage" || v == "--wipe-storage" {
			return []Violation{{
				KataID: "ZC1749",
				Message: "`virsh undefine ... " + v + "` deletes every disk image the " +
					"domain references — no soft-delete, no recycle bin. Back up " +
					"first (`qemu-img convert`), `undefine` without the flag, then " +
					"`virsh vol-delete` deliberately.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

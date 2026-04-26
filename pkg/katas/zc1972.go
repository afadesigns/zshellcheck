// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1972",
		Title:    "Error on `dmsetup remove_all` / `dmsetup remove -f` — tears down live LVM/LUKS/multipath mappings",
		Severity: SeverityError,
		Description: "`dmsetup remove_all` iterates every device-mapper node on the host — " +
			"LVM logical volumes, LUKS containers, multipath aggregates, `cryptsetup` " +
			"mappings — and asks the kernel to drop each one. `dmsetup remove --force " +
			"$NAME` targets a single mapping but still evicts it with in-flight I/O. " +
			"When any of those devices is mounted or backing a running VM, new I/O to " +
			"it returns `ENXIO`, `fsck` is no longer possible, and LVM metadata needs a " +
			"cold reboot to reappear. Use `dmsetup remove $NAME` without `--force` " +
			"after `umount`/`vgchange -an`/`cryptsetup close`, and never `remove_all` " +
			"on a host you care about.",
		Check: checkZC1972,
	})
}

func checkZC1972(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "dmsetup" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	sub := cmd.Arguments[0].String()
	if sub == "remove_all" {
		return zc1972Hit(cmd, "dmsetup remove_all")
	}
	if sub == "remove" {
		for _, arg := range cmd.Arguments[1:] {
			v := arg.String()
			if v == "-f" || v == "--force" {
				return zc1972Hit(cmd, "dmsetup remove -f")
			}
		}
	}
	return nil
}

func zc1972Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1972",
		Message: "`" + form + "` drops LVM/LUKS/multipath mappings while still in " +
			"use — in-flight I/O returns `ENXIO`, metadata needs a reboot. `umount` " +
			"+ `vgchange -an` / `cryptsetup close` first, then `dmsetup remove`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

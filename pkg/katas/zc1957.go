// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1957",
		Title:    "Warn on `lvchange -an` / `vgchange -an` — deactivates a live LV/VG, risks mounted-fs corruption",
		Severity: SeverityWarning,
		Description: "`lvchange -an VG/LV` (and `vgchange -an VG` for the whole group) deactivates " +
			"a logical volume by removing its device-mapper entry. If the LV is mounted, " +
			"writes that the kernel has buffered but not yet flushed may be lost, and any " +
			"process holding an open fd on the filesystem gets EIO on the next syscall. " +
			"`umount` the mount first, stop any service keeping files open, verify with " +
			"`lsof` / `fuser`, and only then `lvchange -an`. For a scripted teardown, prefer " +
			"`umount` + `lvremove` with a recovery snapshot in hand.",
		Check: checkZC1957,
	})
}

func checkZC1957(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "lvchange" && ident.Value != "vgchange" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "-an", "-aN":
			return zc1957Hit(cmd, ident.Value+" -an")
		case "--activate=n", "--activate=N":
			return zc1957Hit(cmd, ident.Value+" "+v)
		case "--activate":
			if i+1 < len(cmd.Arguments) {
				next := cmd.Arguments[i+1].String()
				if next == "n" || next == "N" {
					return zc1957Hit(cmd, ident.Value+" --activate n")
				}
			}
		}
	}
	return nil
}

func zc1957Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1957",
		Message: "`" + form + "` deactivates the LV/VG — unflushed writes on a mounted " +
			"fs may be lost, open fds see EIO. Umount and stop holders first, verify with " +
			"`lsof`/`fuser`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

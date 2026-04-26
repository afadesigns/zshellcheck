// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1526",
		Title:    "Error on `wipefs -a` / `wipefs -af` — erases filesystem signatures (unrecoverable)",
		Severity: SeverityError,
		Description: "`wipefs -a` overwrites every filesystem, partition table, and RAID signature " +
			"it finds on the target. Unlike `rm`, there is no retention anywhere — the only " +
			"recovery path is a disk image backup taken beforehand. If the target variable is " +
			"wrong (typo, empty, resolves to the wrong `/dev/sdX`), this bricks the disk. " +
			"Always run with `--no-act` first or prefer `sgdisk --zap-all` for partition-table " +
			"scope.",
		Check: checkZC1526,
	})
}

func checkZC1526(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "wipefs" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-a" || v == "-af" || v == "-fa" || v == "--all" {
			return []Violation{{
				KataID: "ZC1526",
				Message: "`wipefs -a` erases every filesystem signature — unrecoverable. Run " +
					"with `--no-act` first, or use `sgdisk --zap-all` for scoped deletion.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

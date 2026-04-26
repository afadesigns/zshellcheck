// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1735",
		Title:    "Error on `efibootmgr -B` — deletes UEFI boot entry, may brick boot",
		Severity: SeverityError,
		Description: "`efibootmgr -B` deletes the currently-selected UEFI boot entry; combined " +
			"with `-b BOOTNUM` it removes the specific entry instead. If that entry was " +
			"the only viable bootloader (or the firmware's removable-media fallback is " +
			"not present), the next reboot drops into the UEFI shell or picks an " +
			"unexpected device — recovery needs console access. Run `efibootmgr -v` first " +
			"to inspect `BootOrder`, ensure a fallback (`/EFI/BOOT/BOOTX64.EFI`) is in " +
			"place, and prefer `efibootmgr -o NEW,ORDER` to demote rather than delete.",
		Check: checkZC1735,
	})
}

func checkZC1735(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "efibootmgr" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-B" {
			return []Violation{{
				KataID: "ZC1735",
				Message: "`efibootmgr -B` deletes a UEFI boot entry — wrong BOOTNUM (or " +
					"missing fallback) leaves the box at the UEFI shell on next reboot. " +
					"Inspect `efibootmgr -v` first; demote via `-o NEW,ORDER` instead " +
					"of deleting.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

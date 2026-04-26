// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1772EraseFlags = map[string]bool{
	"--security-erase":          true,
	"--security-erase-enhanced": true,
	"--security-set-pass":       true,
	"--security-unlock":         true,
	"--security-disable":        true,
	"--security-freeze":         true,
	"--trim-sector-ranges":      true,
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1772",
		Title:    "Error on `hdparm --security-erase` / `--trim-sector-ranges` — ATA-level data destruction",
		Severity: SeverityError,
		Description: "`hdparm --security-erase PASS $DISK` issues the ATA `SECURITY ERASE UNIT` " +
			"command: the drive firmware wipes every block, ignoring filesystem or partition " +
			"boundaries, and the operation cannot be interrupted or rolled back. " +
			"`--security-erase-enhanced` is the same but also clears reallocated sectors, and " +
			"`--trim-sector-ranges` discards the listed LBAs on any TRIM-capable device. " +
			"`--security-set-pass`, `--security-disable`, `--security-unlock`, and " +
			"`--security-freeze` alter the drive-level password state and, if misused in a " +
			"script, lock the device out of future access. Keep these calls behind a guarded " +
			"runbook with the exact disk pinned by `/dev/disk/by-id/…` and the password stored " +
			"outside argv.",
		Check: checkZC1772,
	})
}

func checkZC1772(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "hdparm" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if zc1772EraseFlags[v] {
			return zc1772Hit(cmd, v)
		}
	}
	return nil
}

func zc1772Hit(cmd *ast.SimpleCommand, flag string) []Violation {
	line, col := FlagArgPosition(cmd, zc1772EraseFlags)
	return []Violation{{
		KataID: "ZC1772",
		Message: "`hdparm " + flag + "` issues an ATA-level operation that ignores " +
			"filesystems and cannot be rolled back. Pin the disk by " +
			"`/dev/disk/by-id/…`, keep it behind a runbook, keep the password out of argv.",
		Line:   line,
		Column: col,
		Level:  SeverityError,
	}}
}

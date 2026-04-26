// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1692",
		Title:    "Error on `kexec -e` — jumps into a new kernel without reboot, no audit trail",
		Severity: SeverityError,
		Description: "`kexec -e` transfers control to whatever kernel image is currently " +
			"loaded via `kexec -l` — there is no firmware reboot, no init re-run, no " +
			"chance for PAM / auditd / systemd hooks to record the transition. Malware " +
			"uses it to pivot into a rootkit kernel while the audit log shows no reboot. " +
			"If the intent is a fast reboot, prefer `systemctl kexec` (writes a wtmp entry " +
			"and flushes filesystems), or just `reboot` / `systemctl reboot` and take the " +
			"firmware cost for the audit trail.",
		Check: checkZC1692,
	})
}

func checkZC1692(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "kexec" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-e" {
			return []Violation{{
				KataID: "ZC1692",
				Message: "`kexec -e` jumps to a preloaded kernel without firmware reboot " +
					"— wtmp / auditd see nothing. Use `systemctl kexec` or a real " +
					"`systemctl reboot` to keep the audit trail.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

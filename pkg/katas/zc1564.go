// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1564",
		Title:    "Warn on `date -s` / `timedatectl set-time` — manual clock change breaks TLS / cron",
		Severity: SeverityWarning,
		Description: "Setting the system clock by hand (`date -s`, `timedatectl set-time`, " +
			"`hwclock --set`) moves wall-clock time enough to invalidate short-lived TLS " +
			"certificates, reset `cron`'s missed-job catch-up, and confuse `systemd.timer` " +
			"units that depend on monotonic math. Use `systemd-timesyncd` / `chrony` / `ntpd` " +
			"for routine correction; reserve manual set for first-boot bootstrap or air-gapped " +
			"recovery and document the action.",
		Check: checkZC1564,
	})
}

func checkZC1564(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value == "date" {
		for _, arg := range cmd.Arguments {
			if arg.String() == "-s" || arg.String() == "--set" {
				return zc1564Violation(cmd, "date -s")
			}
		}
	}
	if ident.Value == "timedatectl" && len(cmd.Arguments) >= 1 &&
		cmd.Arguments[0].String() == "set-time" {
		return zc1564Violation(cmd, "timedatectl set-time")
	}
	if ident.Value == "hwclock" {
		for _, arg := range cmd.Arguments {
			v := arg.String()
			if v == "--set" || v == "-w" || v == "--systohc" {
				return zc1564Violation(cmd, "hwclock "+v)
			}
		}
	}
	return nil
}

func zc1564Violation(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1564",
		Message: "`" + what + "` sets the wall clock manually — breaks TLS certs, cron " +
			"catch-up, and systemd timer math. Use timesyncd/chrony/ntpd.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

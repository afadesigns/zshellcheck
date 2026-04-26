// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1467",
		Title:    "Warn on `sysctl -w kernel.core_pattern=|...` / `kernel.modprobe=...` (kernel hijack)",
		Severity: SeverityError,
		Description: "Writing `kernel.core_pattern` to a pipe handler or `kernel.modprobe` to a " +
			"user-writable path is a textbook privilege-escalation trick: the next crashing " +
			"setuid process (or the next auto-load of an absent module) executes the supplied " +
			"binary as root. Keep `core_pattern` set to `core` or `systemd-coredump` and leave " +
			"`kernel.modprobe` at the distro default (`/sbin/modprobe`).",
		Check: checkZC1467,
	})
}

func checkZC1467(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "sysctl" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := stripOuterQuotes(arg.String())
		// Accept both `key=value` and `-w key=value` — `-w` shows up as its own arg.
		k, val, ok := strings.Cut(v, "=")
		if !ok {
			continue
		}
		k = strings.TrimSpace(k)
		val = stripOuterQuotes(val)
		if k == "kernel.core_pattern" && strings.HasPrefix(val, "|") {
			return zc1467Violation(cmd, "kernel.core_pattern pipe handler")
		}
		if k == "kernel.modprobe" && val != "" && val != "/sbin/modprobe" {
			return zc1467Violation(cmd, "kernel.modprobe override")
		}
	}
	return nil
}

func stripOuterQuotes(s string) string {
	if len(s) >= 2 {
		first, last := s[0], s[len(s)-1]
		if (first == '\'' && last == '\'') || (first == '"' && last == '"') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func zc1467Violation(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1467",
		Message: "Kernel hijack vector (" + what + ") — next crash / module load runs " +
			"attacker-supplied binary as root.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

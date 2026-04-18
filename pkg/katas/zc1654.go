package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1654",
		Title:    "Warn on `sysctl -p /tmp/...` — loading kernel tunables from attacker-writable path",
		Severity: SeverityWarning,
		Description: "`sysctl -p PATH` reads `key=value` lines from PATH and applies them as " +
			"kernel tunables. A PATH under `/tmp/` or `/var/tmp/` is world-traversable; a " +
			"concurrent local user can substitute the file between write and read, " +
			"injecting `kernel.core_pattern=|/tmp/evil`, `kernel.modprobe=/tmp/evil`, or " +
			"disabling hardening knobs (`kernel.kptr_restrict=0`, `kernel.yama.ptrace_scope=" +
			"0`). Keep sysctl configs under `/etc/sysctl.d/` with root ownership.",
		Check: checkZC1654,
	})
}

func checkZC1654(node ast.Node) []Violation {
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

	for i, arg := range cmd.Arguments {
		if arg.String() != "-p" {
			continue
		}
		if i+1 >= len(cmd.Arguments) {
			continue
		}
		next := cmd.Arguments[i+1].String()
		if strings.HasPrefix(next, "/tmp/") || strings.HasPrefix(next, "/var/tmp/") {
			return []Violation{{
				KataID: "ZC1654",
				Message: "`sysctl -p " + next + "` reads tunables from a world-traversable " +
					"path — a concurrent local user can substitute the file. Keep configs " +
					"under `/etc/sysctl.d/`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1710",
		Title:    "Error on `journalctl --vacuum-size=1` / `--vacuum-time=1s` — journal-wipe pattern",
		Severity: SeverityError,
		Description: "`journalctl --vacuum-size=1` (down to 1 byte / 1K), `--vacuum-time=1s` " +
			"(retain only the last second), or `--vacuum-files=1` (keep one journal file) " +
			"effectively flushes the entire systemd journal. The classic shape after a " +
			"compromise — clear the audit trail before re-enabling logging. Real retention " +
			"belongs in `/etc/systemd/journald.conf` (`SystemMaxUse=`, `MaxRetentionSec=`), " +
			"not in an ad-hoc one-shot. If you genuinely need to bound disk use, set the " +
			"limit to a meaningful value (`--vacuum-time=2weeks`, `--vacuum-size=200M`).",
		Check: checkZC1710,
	})
}

var zc1710VacuumPrefixes = []string{
	"--vacuum-size=",
	"--vacuum-time=",
	"--vacuum-files=",
}

func checkZC1710(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "journalctl" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		for _, prefix := range zc1710VacuumPrefixes {
			if !strings.HasPrefix(v, prefix) {
				continue
			}
			val := strings.TrimPrefix(v, prefix)
			if !zc1710Aggressive(val) {
				continue
			}
			return []Violation{{
				KataID: "ZC1710",
				Message: "`journalctl " + v + "` flushes the systemd journal — classic " +
					"audit-clear shape. Set retention in `/etc/systemd/journald.conf` " +
					"(`SystemMaxUse=`, `MaxRetentionSec=`) instead.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

// zc1710Aggressive returns true for vacuum values that effectively wipe the
// journal: `1`, `1B`, `1K`, `1KB`, `1s`, `1m`, `0`, etc.
func zc1710Aggressive(val string) bool {
	switch val {
	case "0", "1":
		return true
	}
	low := strings.ToLower(val)
	switch low {
	case "1b", "1k", "1kb", "1kib", "1s", "1m", "1ms", "1µs", "0s", "0m":
		return true
	}
	return false
}

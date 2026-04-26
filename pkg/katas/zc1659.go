// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1659",
		Title:    "Warn on `fuser -k <path>` — kills every process holding the subtree open",
		Severity: SeverityWarning,
		Description: "`fuser -k PATH` sends a signal (SIGKILL by default) to every process that " +
			"has any file under PATH open — not just the one you expected. On `/`, `/var`, " +
			"`/tmp`, or any mount-root this reaches sshd, cron, dbus, and the caller's own " +
			"shell; on a bind-mount it kills workloads that share the host inode. Target " +
			"specific PIDs (`kill $(pidof app)`) or ports (`fuser -k PORT/tcp`), or use " +
			"`systemctl stop UNIT` for services. `fuser -k` against a filesystem path is " +
			"blast-radius that the caller rarely owns.",
		Check: checkZC1659,
	})
}

func checkZC1659(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "fuser" {
		return nil
	}

	hasKill := false
	pathTarget := ""
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.HasPrefix(v, "-") && !strings.HasPrefix(v, "--") {
			if strings.ContainsRune(strings.TrimPrefix(v, "-"), 'k') {
				hasKill = true
			}
			continue
		}
		if strings.HasPrefix(v, "/") {
			if strings.HasSuffix(v, "/tcp") || strings.HasSuffix(v, "/udp") ||
				strings.HasSuffix(v, "/sctp") {
				continue
			}
			pathTarget = v
		}
	}

	if !hasKill || pathTarget == "" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1659",
		Message: "`fuser -k " + pathTarget + "` signals every process with a file open " +
			"anywhere under the path — use PID / port targets or `systemctl stop` for " +
			"services.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

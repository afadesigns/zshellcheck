// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1627",
		Title:    "Warn on `crontab /tmp/FILE` — attacker-writable path installed as a user's cron",
		Severity: SeverityWarning,
		Description: "`crontab PATH` replaces the user's cron with whatever PATH currently " +
			"contains. A path under `/tmp/` or `/var/tmp/` is world-traversable; a concurrent " +
			"local user can replace the file between the moment the script writes it and the " +
			"moment `crontab` reads it, substituting their own cron rules. Keep the staging " +
			"file in a 0700-scoped directory (e.g. `$XDG_RUNTIME_DIR/` or `mktemp -d`), or " +
			"pipe the content via `crontab -` after generating it in-memory.",
		Check: checkZC1627,
	})
}

func checkZC1627(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "crontab" {
		return nil
	}

	var skipNext bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if skipNext {
			skipNext = false
			continue
		}
		if v == "-u" || v == "-s" {
			skipNext = true
			continue
		}
		if strings.HasPrefix(v, "-") {
			continue
		}
		if strings.HasPrefix(v, "/tmp/") || strings.HasPrefix(v, "/var/tmp/") {
			return []Violation{{
				KataID: "ZC1627",
				Message: "`crontab " + v + "` reads cron rules from a world-traversable " +
					"path — a concurrent local user can substitute the file between write " +
					"and read. Stage the file in `$XDG_RUNTIME_DIR/` or `mktemp -d`, or " +
					"pipe via `crontab -`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

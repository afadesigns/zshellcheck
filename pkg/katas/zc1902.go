// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1902SensitiveTargets = []string{
	"/var/log/",
	"/var/log/audit/",
	"/var/log/wtmp",
	"/var/log/btmp",
	"/var/log/lastlog",
	"/var/log/secure",
	"/var/log/auth.log",
	"/var/log/syslog",
	"/var/log/messages",
	"/.bash_history",
	"/.zsh_history",
	"/.ash_history",
	"/.python_history",
	"/.mysql_history",
	"/.psql_history",
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1902",
		Title:    "Error on `ln -s /dev/null <logfile>` — silently discards audit or history writes",
		Severity: SeverityError,
		Description: "A symlink from an audit or shell-history path to `/dev/null` turns every " +
			"subsequent append into a no-op — `/var/log/auth.log`, `wtmp`, `~/.bash_history`, " +
			"`~/.zsh_history` all stop recording. This is the textbook way to cover tracks on " +
			"a compromised host and almost never appears in benign automation. If you really " +
			"need to stop a log, disable the writer (rsyslog rule, `set +o history`) or rotate " +
			"with `logrotate` — never redirect into `/dev/null`.",
		Check: checkZC1902,
	})
}

func checkZC1902(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ln" {
		return nil
	}

	var symbolic bool
	var source, target string
	positional := 0
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.HasPrefix(v, "-") {
			if strings.ContainsAny(v, "sS") {
				symbolic = true
			}
			continue
		}
		switch positional {
		case 0:
			source = v
		case 1:
			target = v
		}
		positional++
	}
	if !symbolic || source != "/dev/null" {
		return nil
	}
	if !zc1902IsSensitive(target) {
		return nil
	}

	return []Violation{{
		KataID: "ZC1902",
		Message: "`ln -s /dev/null " + target + "` redirects every write to the " +
			"bit-bucket — audit / history entries vanish silently. If the log must " +
			"stop, disable the writer or rotate with `logrotate`, never symlink to `/dev/null`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

func zc1902IsSensitive(target string) bool {
	if target == "" {
		return false
	}
	for _, suffix := range zc1902SensitiveTargets {
		if strings.HasSuffix(target, suffix) || strings.Contains(target, suffix) {
			return true
		}
	}
	return false
}

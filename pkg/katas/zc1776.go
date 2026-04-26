// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

// Schemes that commonly embed credentials in a connection URI and are
// passed to a CLI client that keeps the URI in argv.
var zc1776CredSchemes = []string{
	"postgres://",
	"postgresql://",
	"mysql://",
	"mariadb://",
	"mongodb://",
	"mongodb+srv://",
	"redis://",
	"rediss://",
	"amqp://",
	"amqps://",
	"kafka://",
	"nats://",
	"clickhouse://",
	"cockroachdb://",
	"db2://",
	"mssql://",
	"oracle://",
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1776",
		Title:    "Error on `psql postgresql://user:secret@host/db` — password in argv via connection URI",
		Severity: SeverityError,
		Description: "Database and message-broker CLIs accept a single connection URI " +
			"(`postgresql://`, `mysql://`, `mongodb://`, `redis://`, `amqp://`, `kafka://`, " +
			"and friends). When the URI embeds a password — `scheme://user:secret@host/db` — " +
			"the secret lands in argv, visible to every user via `ps`, `/proc/PID/cmdline`, " +
			"process accounting, and audit trails, and it often survives in shell history. " +
			"Keep the password out of argv: use the client's password-file / `.pgpass` / " +
			"`PGPASSWORD` / `REDISCLI_AUTH` equivalent, or interpolate the URI from an " +
			"environment variable so the secret is not on the command line that other users " +
			"can see.",
		Check: checkZC1776,
	})
}

func checkZC1776(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if _, ok := cmd.Name.(*ast.Identifier); !ok {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		v = strings.Trim(v, "\"'")
		if leak, scheme := zc1776UriHasPassword(v); leak {
			return []Violation{{
				KataID: "ZC1776",
				Message: "`" + scheme + "user:SECRET@…` in argv puts the password in `ps` / " +
					"`/proc/PID/cmdline` / history. Use a password file (`~/.pgpass`, " +
					"`~/.my.cnf`), `PGPASSWORD` / `REDISCLI_AUTH` env var, or build the URI " +
					"from a secret variable.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

func zc1776UriHasPassword(v string) (bool, string) {
	for _, scheme := range zc1776CredSchemes {
		if !strings.HasPrefix(v, scheme) {
			continue
		}
		rest := v[len(scheme):]
		at := strings.Index(rest, "@")
		if at <= 0 {
			return false, scheme
		}
		userinfo := rest[:at]
		colon := strings.Index(userinfo, ":")
		if colon <= 0 || colon == len(userinfo)-1 {
			// No password segment, or empty password.
			return false, scheme
		}
		return true, scheme
	}
	return false, ""
}

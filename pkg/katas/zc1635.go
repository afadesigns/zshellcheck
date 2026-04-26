// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1635",
		Title:    "Error on `mysql -pSECRET` / `--password=SECRET` — password in process list",
		Severity: SeverityError,
		Description: "MySQL / MariaDB clients accept the password concatenated with the `-p` " +
			"flag (`-pSECRET`) or via `--password=SECRET`. Both forms put the secret in argv " +
			"— visible in `ps`, `/proc/<pid>/cmdline`, shell history, and audit logs for every " +
			"local user who can list processes. Use `-p` with no argument for an interactive " +
			"prompt, `--login-path` for the credentials helper file, or a `~/.my.cnf` with " +
			"`0600` perms and `[client] password=...` so the client reads it at startup.",
		Check: checkZC1635,
	})
}

func checkZC1635(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	switch ident.Value {
	case "mysql", "mysqldump", "mysqladmin", "mariadb", "mariadb-dump", "mariadb-admin":
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.HasPrefix(v, "-p") && len(v) > 2 {
			return []Violation{{
				KataID: "ZC1635",
				Message: "`" + ident.Value + " " + v + "` puts the MySQL password in argv. " +
					"Use `-p` with no arg (prompt), `--login-path`, or a 0600 `~/.my.cnf`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
		if strings.HasPrefix(v, "--password=") {
			return []Violation{{
				KataID: "ZC1635",
				Message: "`" + ident.Value + " " + v + "` puts the MySQL password in argv. " +
					"Use `-p` with no arg (prompt), `--login-path`, or a 0600 `~/.my.cnf`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

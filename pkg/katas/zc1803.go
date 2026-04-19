package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1803MySQLClients = map[string]bool{
	"mysql":         true,
	"mysqldump":     true,
	"mysqladmin":    true,
	"mariadb":       true,
	"mariadb-dump":  true,
	"mariadb-admin": true,
}

var zc1803PgClients = map[string]bool{
	"psql":       true,
	"pg_dump":    true,
	"pgbench":    true,
	"pg_restore": true,
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1803",
		Title:    "Error on `mysql --skip-ssl` / `psql sslmode=disable` — plaintext credentials on the wire",
		Severity: SeverityError,
		Description: "Disabling TLS on a MySQL or PostgreSQL client pushes the login handshake " +
			"(including the password or auth challenge) and every subsequent query and " +
			"result over plaintext TCP. Anyone in the network path — the cloud VPC, the " +
			"office LAN, a compromised router — can sniff or modify the stream. The flags " +
			"vary (`--skip-ssl`, `--ssl=0`, `--ssl-mode=DISABLED` for MySQL / MariaDB; " +
			"`sslmode=disable` in the connection URI or `PGSSLMODE=disable` env var for " +
			"PostgreSQL) but the effect is the same. Prefer `--ssl-mode=VERIFY_IDENTITY` " +
			"(MySQL 8+) and `sslmode=verify-full` (psql) with a pinned CA bundle.",
		Check: checkZC1803,
	})
}

func checkZC1803(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	// Parser caveat: `mysql --skip-ssl …` mangles name to `skip-ssl`.
	if ident.Value == "skip-ssl" {
		return zc1803Hit(cmd, "mysql", "--skip-ssl")
	}

	if zc1803MySQLClients[ident.Value] {
		for _, arg := range cmd.Arguments {
			raw := strings.Trim(arg.String(), "\"'")
			v := strings.ToLower(raw)
			if v == "--skip-ssl" || v == "--ssl=0" || v == "--ssl=false" ||
				v == "--ssl-mode=disabled" || v == "--ssl-mode=disable" {
				return zc1803Hit(cmd, ident.Value, raw)
			}
		}
	}

	if zc1803PgClients[ident.Value] {
		for _, arg := range cmd.Arguments {
			raw := strings.Trim(arg.String(), "\"'")
			if strings.Contains(strings.ToLower(raw), "sslmode=disable") {
				return zc1803Hit(cmd, ident.Value, raw)
			}
		}
	}
	return nil
}

func zc1803Hit(cmd *ast.SimpleCommand, tool, flag string) []Violation {
	return []Violation{{
		KataID: "ZC1803",
		Message: "`" + tool + " " + flag + "` disables TLS — login handshake and " +
			"queries travel in plaintext. Use `--ssl-mode=VERIFY_IDENTITY` (MySQL) / " +
			"`sslmode=verify-full` (psql) with a pinned CA.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1445",
		Title:    "Dangerous: `dropdb` / `mysqladmin drop` — deletes a database",
		Severity: SeverityError,
		Description: "`dropdb NAME` removes a PostgreSQL database including all data and " +
			"schemas. `mysqladmin drop NAME` does the same for MySQL. Always `pg_dump` / " +
			"`mysqldump` first and consider requiring `-i`/`-y`-less forms so operators must " +
			"type confirmation.",
		Check: checkZC1445,
	})
}

func checkZC1445(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "dropdb":
		return []Violation{{
			KataID:  "ZC1445",
			Message: "`dropdb` removes a PostgreSQL database. Verify target and backup first (`pg_dump`).",
			Line:    cmd.Token.Line,
			Column:  cmd.Token.Column,
			Level:   SeverityError,
		}}
	case "mysqladmin":
		for _, arg := range cmd.Arguments {
			if arg.String() == "drop" {
				return []Violation{{
					KataID:  "ZC1445",
					Message: "`mysqladmin drop` removes a MySQL database. Verify target and backup first (`mysqldump`).",
					Line:    cmd.Token.Line,
					Column:  cmd.Token.Column,
					Level:   SeverityError,
				}}
			}
		}
	}

	return nil
}

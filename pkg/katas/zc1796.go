// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1796",
		Title:    "Warn on `pg_restore --clean` / `-c` — drops existing DB objects before recreating",
		Severity: SeverityWarning,
		Description: "`pg_restore -c` (also `--clean`) issues `DROP` for every table, index, " +
			"function, and sequence in the target database before recreating them from the " +
			"archive. If the backup is stale, incomplete, or points at the wrong database, " +
			"the destination loses any object that isn't in the dump — including data added " +
			"after the backup ran. Restore into a fresh empty database (`createdb new && " +
			"pg_restore -d new`) or snapshot the target (`pg_dump -Fc > pre.dump`) before " +
			"running `--clean`, and never pair it with `--if-exists` on a live production DB " +
			"without a tested rollback path.",
		Check: checkZC1796,
	})
}

func checkZC1796(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	// Parser caveat: `pg_restore --clean …` mangles to name=`clean`.
	// To avoid false positives on unrelated commands with `clean` as name,
	// require another pg_restore-ish argument to be present.
	if ident.Value == "clean" {
		if zc1796HasPgArg(cmd) {
			return zc1796Hit(cmd, "pg_restore --clean")
		}
		return nil
	}

	if ident.Value != "pg_restore" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-c" || v == "--clean" {
			return zc1796Hit(cmd, "pg_restore "+v)
		}
	}
	return nil
}

func zc1796HasPgArg(cmd *ast.SimpleCommand) bool {
	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "-d", "--dbname", "-F", "--format", "-U", "--username",
			"--if-exists", "--no-owner", "--no-acl":
			return true
		}
	}
	return false
}

func zc1796Hit(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1796",
		Message: "`" + what + "` drops every object in the target DB before recreating " +
			"from the archive — stale or wrong-target dump silently loses data. Restore " +
			"into a fresh DB (`createdb new && pg_restore -d new`), or snapshot first.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

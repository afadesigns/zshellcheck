// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1738",
		Title:    "Error on `aws rds delete-db-instance --skip-final-snapshot` — DB destroyed unrecoverable",
		Severity: SeverityError,
		Description: "RDS keeps a final snapshot when an instance or cluster is deleted — the only " +
			"path back from a typo'd identifier or wrong account. `--skip-final-snapshot` " +
			"opts out of that snapshot, so the database is gone the moment the API call " +
			"returns; same applies to `aws rds delete-db-cluster --skip-final-snapshot`. " +
			"Drop the flag (or pass `--final-db-snapshot-identifier <name>` so the snapshot " +
			"name is explicit) and verify the snapshot lands before reusing the identifier.",
		Check: checkZC1738,
	})
}

func checkZC1738(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "aws" {
		return nil
	}
	if len(cmd.Arguments) < 3 {
		return nil
	}
	if cmd.Arguments[0].String() != "rds" {
		return nil
	}
	sub := cmd.Arguments[1].String()
	if sub != "delete-db-instance" && sub != "delete-db-cluster" {
		return nil
	}

	for _, arg := range cmd.Arguments[2:] {
		if arg.String() == "--skip-final-snapshot" {
			return []Violation{{
				KataID: "ZC1738",
				Message: "`aws rds " + sub + " --skip-final-snapshot` deletes the database " +
					"with no recovery snapshot. Drop the flag or pass `--final-db-" +
					"snapshot-identifier <name>` so the snapshot is explicit and " +
					"verifiable.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

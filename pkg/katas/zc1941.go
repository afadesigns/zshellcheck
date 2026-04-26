// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1941",
		Title:    "Error on `restic init --insecure-no-password` — creates an unencrypted backup repository",
		Severity: SeverityError,
		Description: "`restic init --insecure-no-password` creates a repo whose data chunks are " +
			"reachable without a key. Every later `backup` and `restore` round-trips " +
			"plaintext blocks to the storage backend, so any operator with read access to the " +
			"bucket / NFS share / SFTP directory can assemble the backed-up filesystem — " +
			"including shell history, SSH keys, and database dumps. Pass a real passphrase via " +
			"`--password-file` (mode `0400`, readable only by the backup user) or " +
			"`--password-command`, and never use the `--insecure-*` family outside a local " +
			"test repo.",
		Check: checkZC1941,
	})
}

func checkZC1941(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	// Parser caveat: `restic --insecure-no-password …` mangles the command
	// name to `insecure-no-password`.
	if ident.Value == "insecure-no-password" {
		return zc1941Hit(cmd)
	}
	if ident.Value != "restic" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		if arg.String() == "--insecure-no-password" {
			return zc1941Hit(cmd)
		}
	}
	return nil
}

func zc1941Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1941",
		Message: "`restic --insecure-no-password` creates an unencrypted repo — every " +
			"operator with read access to the backend can reassemble the backed-up " +
			"filesystem. Use `--password-file` / `--password-command` with a real passphrase.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1633",
		Title:    "Error on `gpg --passphrase SECRET` — passphrase on cmdline",
		Severity: SeverityError,
		Description: "`gpg --passphrase VALUE` passes the key passphrase as an argv element. " +
			"Visible in `ps`, `/proc/<pid>/cmdline`, shell history, and audit logs for every " +
			"local user who can list processes. Use `--passphrase-file PATH` (reads the first " +
			"line of PATH), `--passphrase-fd N` (reads from file descriptor N), or " +
			"`--pinentry-mode=loopback` with the passphrase piped on stdin. Pair with " +
			"`--batch` for non-interactive runs.",
		Check: checkZC1633,
	})
}

func checkZC1633(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "gpg" && ident.Value != "gpg2" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "--passphrase" {
			return []Violation{{
				KataID: "ZC1633",
				Message: "`gpg --passphrase` puts the passphrase in argv — visible via `ps`. " +
					"Use `--passphrase-file`, `--passphrase-fd`, or `--pinentry-" +
					"mode=loopback` with the value on stdin.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1614",
		Title:    "Error on `expect` script containing `password` / `passphrase`",
		Severity: SeverityError,
		Description: "`expect -c '... password ... send \"...\"'` puts the entire scripted " +
			"dialog on the command line. Anything there — including the password or passphrase " +
			"— is visible in `ps`, `/proc/<pid>/cmdline`, shell history, and audit logs. Use " +
			"key-based authentication (SSH keys, GSSAPI) where possible. If password feeding is " +
			"truly unavoidable, read it from a protected file with `spawn -o`, or source it " +
			"from an environment variable the script does not print.",
		Check: checkZC1614,
	})
}

func checkZC1614(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "expect" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		low := strings.ToLower(arg.String())
		if strings.Contains(low, "password") || strings.Contains(low, "passphrase") {
			return []Violation{{
				KataID: "ZC1614",
				Message: "`expect` script contains `password` / `passphrase` — the full " +
					"argv lands in `ps` and audit logs. Switch to key-based auth, or " +
					"read the credential from a protected file the expect script opens.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

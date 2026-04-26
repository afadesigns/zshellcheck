// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1421",
		Title:    "Avoid `chpasswd` / `passwd --stdin` — plaintext passwords in process tree",
		Severity: SeverityError,
		Description: "Passing passwords on stdin to `chpasswd` or `passwd --stdin` exposes the " +
			"plaintext in the process command line or pipeline — visible to `ps`, logs, and " +
			"environment. Use encrypted-hash input (`chpasswd -e`), `usermod -p` with a hash, " +
			"or an IaC tool that handles credentials outside the process tree.",
		Check: checkZC1421,
	})
}

func checkZC1421(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "chpasswd" {
		return nil
	}

	// Any chpasswd invocation without -e (encrypted) is risky.
	hasEncrypted := false
	for _, arg := range cmd.Arguments {
		if arg.String() == "-e" || arg.String() == "--encrypted" {
			hasEncrypted = true
		}
	}
	if !hasEncrypted {
		return []Violation{{
			KataID: "ZC1421",
			Message: "`chpasswd` without `-e`/`--encrypted` accepts plaintext passwords — avoid " +
				"piping cleartext credentials into the process tree. Use a password hash (`-e`) " +
				"or a credentials store.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityError,
		}}
	}

	return nil
}

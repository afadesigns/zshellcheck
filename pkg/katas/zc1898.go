// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1898",
		Title:    "Error on `gpg --export-secret-keys` — private-key material leaks to stdout",
		Severity: SeverityError,
		Description: "`gpg --export-secret-keys KEYID` and `--export-secret-subkeys` write the " +
			"ASCII-armoured private key to stdout. In a script, that stream usually lands " +
			"in a file the operator plans to move off-box — and any misstep (wrong " +
			"`cd`, script-wide stdout captured by CI, tee to a world-readable log, " +
			"piped into a remote unencrypted channel) permanently leaks the key. Backup " +
			"the key interactively on an air-gapped machine; if automation is required, " +
			"write the output to a `umask 077`-protected path and immediately encrypt " +
			"with a second symmetric passphrase.",
		Check: checkZC1898,
	})
}

var zc1898ExportFlags = map[string]bool{
	"--export-secret-keys":    true,
	"--export-secret-subkeys": true,
}

func checkZC1898(node ast.Node) []Violation {
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
		v := arg.String()
		if zc1898ExportFlags[v] {
			line, col := FlagArgPosition(cmd, zc1898ExportFlags)
			return []Violation{{
				KataID: "ZC1898",
				Message: "`gpg " + v + "` writes the private key to stdout — one " +
					"CI-log or wrong-tty redirect leaks it. Back up interactively on an " +
					"air-gapped host, or write to a `umask 077` path and re-encrypt.",
				Line:   line,
				Column: col,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

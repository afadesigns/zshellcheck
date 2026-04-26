// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1840",
		Title:    "Error on `openssl enc -k PASSWORD` — legacy flag embeds secret in argv",
		Severity: SeverityError,
		Description: "`openssl enc -k PASSWORD` (the pre-OpenSSL-3 short form of `-pass " +
			"pass:PASSWORD`) takes the password directly as the next argv element — which " +
			"makes it visible to every `ps` reader, every `/proc/<pid>/cmdline` consumer, " +
			"shell history, and anything that logs command invocations. The same leak " +
			"applies to `openssl rsa`, `openssl pkcs12`, and other subcommands that still " +
			"accept the deprecated `-k` alias. Use `-pass env:VARNAME`, `-pass file:PATH`, " +
			"or `-pass fd:N` (read from an open descriptor) so the secret never rides in " +
			"the process argument vector.",
		Check: checkZC1840,
	})
}

func checkZC1840(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "openssl" {
		return nil
	}
	args := cmd.Arguments
	for i, arg := range args {
		if arg.String() != "-k" {
			continue
		}
		if i+1 >= len(args) {
			continue
		}
		val := args[i+1].String()
		// Ignore empty, flag-looking, or the `-k file:` / `-k env:` style
		// that newer openssl binaries tolerate.
		if val == "" || val[0] == '-' {
			continue
		}
		return []Violation{{
			KataID: "ZC1840",
			Message: "`openssl -k " + val + "` embeds the password in argv — visible " +
				"to `ps`, `/proc/<pid>/cmdline`, and shell history. Use " +
				"`-pass env:VAR`, `-pass file:PATH`, or `-pass fd:N`.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityError,
		}}
	}
	return nil
}

// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1182",
		Title:    "Avoid `nc`/`netcat` for HTTP — use `curl` or `zsh/net/tcp`",
		Severity: SeverityWarning,
		Description: "`nc`/`netcat` for HTTP requests is fragile and lacks TLS support. " +
			"Use `curl` or Zsh `zsh/net/tcp` module for reliable network operations.",
		Check: checkZC1182,
	})
}

func checkZC1182(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "nc" && ident.Value != "netcat" && ident.Value != "ncat" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1182",
		Message: "Avoid `" + ident.Value + "` for network operations in scripts. Use `curl` for HTTP " +
			"or `zmodload zsh/net/tcp` for raw TCP connections with TLS support.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

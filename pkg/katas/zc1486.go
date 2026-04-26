// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1486",
		Title:    "Warn on `curl -2` / `-3` — forces broken SSLv2 / SSLv3",
		Severity: SeverityWarning,
		Description: "`curl -2` (SSLv2) and `-3` (SSLv3) force protocols that are removed from " +
			"every current TLS library. `-2` matches no working server; `-3` leaves you open to " +
			"POODLE. If the remote really needs an old protocol the fix is on the server, not " +
			"the client.",
		Check: checkZC1486,
	})
}

func checkZC1486(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "curl" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-2" || v == "-3" {
			return []Violation{{
				KataID: "ZC1486",
				Message: "`curl " + v + "` forces SSLv2/SSLv3 — removed from modern TLS " +
					"libraries and subject to POODLE. Fix the server instead.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

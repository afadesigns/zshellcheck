// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1485",
		Title:    "Warn on `openssl s_client -ssl3 / -tls1 / -tls1_1` — legacy TLS",
		Severity: SeverityWarning,
		Description: "Forcing SSLv3, TLSv1.0, or TLSv1.1 connects with protocols that have known " +
			"downgrade and bit-flip attacks (POODLE, BEAST). These are disabled by default in " +
			"every maintained OpenSSL build. If the remote only speaks an old protocol, the " +
			"right fix is to update the remote, not downgrade your client.",
		Check: checkZC1485,
	})
}

func checkZC1485(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "openssl" {
		return nil
	}

	if len(cmd.Arguments) == 0 {
		return nil
	}
	sub := cmd.Arguments[0].String()
	if sub != "s_client" && sub != "s_server" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "-ssl2" || v == "-ssl3" || v == "-tls1" || v == "-tls1_1" ||
			v == "-no_tls1_2" || v == "-no_tls1_3" {
			return []Violation{{
				KataID: "ZC1485",
				Message: "`openssl " + sub + " " + v + "` forces a legacy / disabled TLS " +
					"version (downgrade-attack surface). Update the remote instead.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

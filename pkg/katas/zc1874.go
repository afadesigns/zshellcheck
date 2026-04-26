// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1874",
		Title:    "Warn on `sshuttle -r HOST 0/0` — every outbound packet tunneled through the jump host",
		Severity: SeverityWarning,
		Description: "`sshuttle -r user@host 0/0` (or `0.0.0.0/0`, `::/0`) installs a VPN-like " +
			"catch-all route: every TCP connection and DNS lookup on the local machine " +
			"egresses through `user@host`, including traffic to corporate VPN endpoints, " +
			"cloud APIs, and package mirrors that had been explicitly split-tunnel. If the " +
			"jump host is compromised, misconfigured, or simply overloaded, every session " +
			"on the workstation silently degrades or leaks to the wrong peer. Scope the " +
			"subnet list to the networks you actually need (`10.0.0.0/8 172.16.0.0/12 " +
			"192.168.0.0/16`), or prefer `ssh -D` with `--exclude` rules for a single " +
			"browser profile.",
		Check: checkZC1874,
	})
}

func checkZC1874(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "sshuttle" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if zc1874IsDefaultRoute(v) {
			return []Violation{{
				KataID: "ZC1874",
				Message: "`sshuttle ... " + v + "` routes every outbound packet through " +
					"the jump host — a compromise of `user@host` sees the whole " +
					"fleet's traffic. Scope to the subnets you actually need.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

func zc1874IsDefaultRoute(v string) bool {
	switch v {
	case "0/0", "0.0.0.0/0", "::/0":
		return true
	}
	return false
}

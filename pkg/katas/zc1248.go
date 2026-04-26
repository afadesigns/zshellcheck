// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1248",
		Title:    "Prefer `ufw`/`firewalld` over raw `iptables`",
		Severity: SeverityInfo,
		Description: "Raw `iptables` rules are complex and non-persistent by default. " +
			"Use `ufw` (Ubuntu) or `firewalld` (RHEL) for manageable, persistent firewall rules.",
		Check: checkZC1248,
	})
}

func checkZC1248(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "iptables" && ident.Value != "ip6tables" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1248",
		Message: "Prefer `ufw` or `firewalld` over raw `iptables`. " +
			"Firewall frontends provide persistent, manageable rules.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}

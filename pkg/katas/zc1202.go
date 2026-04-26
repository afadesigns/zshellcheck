// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1202",
		Title:    "Avoid `ifconfig` — use `ip` for network configuration",
		Severity: SeverityInfo,
		Description: "`ifconfig` is deprecated on modern Linux. " +
			"Use `ip addr`, `ip link`, or `ip route` from iproute2 for network operations.",
		Check: checkZC1202,
		Fix:   fixZC1202,
	})
}

// fixZC1202 rewrites `ifconfig` to `ip addr` at the command name
// position. `ip addr` is the closest single-token-equivalent iproute2
// invocation; arguments stay untouched and operators/flags must be
// adjusted manually for non-trivial cases.
func fixZC1202(node ast.Node, v Violation, _ []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ifconfig" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("ifconfig"),
		Replace: "ip addr",
	}}
}

func checkZC1202(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ifconfig" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1202",
		Message: "Avoid `ifconfig` — it is deprecated on modern Linux. " +
			"Use `ip addr`, `ip link`, or `ip route` from iproute2.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}

// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1947",
		Title:    "Error on `ip xfrm state flush` / `ip xfrm policy flush` — tears down every IPsec SA and policy",
		Severity: SeverityError,
		Description: "`ip xfrm state flush` removes every IPsec Security Association; " +
			"`ip xfrm policy flush` removes every policy that would have driven them. " +
			"Strongswan, libreswan, FRR, and WireGuard-over-xfrm all lose their tunnels " +
			"instantly — site-to-site VPNs drop, kernel packet paths stop encrypting, and " +
			"peers renegotiate from scratch (with traffic leaking in plaintext during the gap " +
			"on misconfigured hosts). Use `ip xfrm state deleteall src $A dst $B` to scope " +
			"the change to a single tunnel, and pair flushes with a maintenance window.",
		Check: checkZC1947,
	})
}

func checkZC1947(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ip" {
		return nil
	}
	if len(cmd.Arguments) < 3 {
		return nil
	}
	if cmd.Arguments[0].String() != "xfrm" {
		return nil
	}
	tbl := cmd.Arguments[1].String()
	if tbl != "state" && tbl != "policy" {
		return nil
	}
	if cmd.Arguments[2].String() != "flush" {
		return nil
	}
	return []Violation{{
		KataID: "ZC1947",
		Message: "`ip xfrm " + tbl + " flush` tears down every IPsec SA/policy — " +
			"VPN tunnels drop, kernel stops encrypting, plaintext may leak during renegotiation. " +
			"Scope via `ip xfrm " + tbl + " deleteall src $A dst $B`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

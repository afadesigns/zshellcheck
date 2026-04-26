// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1536",
		Title:    "Warn on `iptables -j DNAT` / `-j REDIRECT` — rewrites traffic destination",
		Severity: SeverityWarning,
		Description: "`-j DNAT` and `-j REDIRECT` in an iptables rule rewrite the destination " +
			"address/port of matching packets. That is how you transparently proxy, but also " +
			"how you silently redirect a victim's connections to an attacker-controlled port. " +
			"Scripts that touch NAT rules should be carefully reviewed; prefer declarative " +
			"network config (nftables ruleset, NetworkManager connection, firewalld service) " +
			"and store rule provenance.",
		Check: checkZC1536,
	})
}

func checkZC1536(node ast.Node) []Violation {
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

	args := make([]string, 0, len(cmd.Arguments))
	for _, a := range cmd.Arguments {
		args = append(args, a.String())
	}

	// Must see an add/insert verb to avoid flagging -L listings.
	var hasAdd bool
	for _, a := range args {
		if a == "-A" || a == "-I" || a == "-R" || a == "--append" ||
			a == "--insert" || a == "--replace" {
			hasAdd = true
			break
		}
	}
	if !hasAdd {
		return nil
	}

	for i, a := range args {
		if a == "-j" && i+1 < len(args) {
			tgt := args[i+1]
			if tgt == "DNAT" || tgt == "REDIRECT" || tgt == "NETMAP" {
				return []Violation{{
					KataID: "ZC1536",
					Message: "`iptables -j " + tgt + "` rewrites packet destination — " +
						"silent redirect surface. Use declarative nftables/firewalld config.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
	}
	return nil
}

// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1439",
		Title:    "Enabling IP forwarding in a script — document firewall posture",
		Severity: SeverityWarning,
		Description: "Setting `net.ipv4.ip_forward=1` (or `-w`-ing a sysctl to the same effect) " +
			"turns the host into a router. Without matching iptables/nftables rules this can " +
			"silently expose services between interfaces. If intentional (VPN, container host, " +
			"NAT gateway), pair with explicit firewall rules and persist via `/etc/sysctl.d/`.",
		Check: checkZC1439,
	})
}

func checkZC1439(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "sysctl" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.Contains(v, "ip_forward=1") ||
			strings.Contains(v, "forwarding=1") ||
			strings.Contains(v, "ip_forward =1") {
			return []Violation{{
				KataID: "ZC1439",
				Message: "Enabling `ip_forward` turns the host into a router. Verify firewall " +
					"posture (iptables/nftables) and persist the setting in `/etc/sysctl.d/`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}

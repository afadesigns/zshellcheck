// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1488",
		Title:    "Warn on `ssh -R 0.0.0.0:...` / `*:...` — reverse tunnel bound to all interfaces",
		Severity: SeverityWarning,
		Description: "The default for `ssh -R` binds the remote listener to `localhost`. Pointing " +
			"it at `0.0.0.0` or `*` (or an explicit public IP) exposes the forwarded port to the " +
			"whole network, including anything else that has reached the jump host. For " +
			"persistent ops tunnels, pin the bind address to a specific private interface and " +
			"require `GatewayPorts clientspecified` server-side.",
		Check: checkZC1488,
	})
}

func checkZC1488(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ssh" && ident.Value != "autossh" {
		return nil
	}

	var prevForward bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevForward {
			prevForward = false
			if strings.HasPrefix(v, "0.0.0.0:") || strings.HasPrefix(v, "*:") ||
				strings.HasPrefix(v, "::") {
				return []Violation{{
					KataID: "ZC1488",
					Message: "SSH reverse tunnel bound to all interfaces (`" + v + "`) — " +
						"forwarded port reachable from any network. Bind to a specific IP.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
		if v == "-R" || v == "-L" || v == "-D" {
			prevForward = true
		}
	}
	return nil
}

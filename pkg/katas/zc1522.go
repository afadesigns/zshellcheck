// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1522",
		Title:    "Warn on `ip route add default` / `route add default` — changes default gateway",
		Severity: SeverityWarning,
		Description: "Setting a new default route in a script silently redirects every non-local " +
			"packet through the specified gateway. That is exactly the knob an attacker turns " +
			"to MITM a whole host after a foothold, and it is also a common accidental foot- " +
			"gun in CI runners (gateway in the runner network ≠ gateway in production). Use " +
			"NetworkManager / systemd-networkd config files for persistent routes, and " +
			"document any runtime change with a comment explaining why.",
		Check: checkZC1522,
	})
}

func checkZC1522(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	args := make([]string, 0, len(cmd.Arguments))
	for _, a := range cmd.Arguments {
		args = append(args, a.String())
	}

	// `ip route add default ...`
	if ident.Value == "ip" && len(args) >= 3 &&
		args[0] == "route" && args[1] == "add" && args[2] == "default" {
		return zc1522Violation(cmd, "ip route add default")
	}
	// `route add default ...`
	if ident.Value == "route" && len(args) >= 2 &&
		args[0] == "add" && args[1] == "default" {
		return zc1522Violation(cmd, "route add default")
	}
	return nil
}

func zc1522Violation(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1522",
		Message: "`" + what + "` silently reroutes every non-local packet through the new " +
			"gateway — MITM surface or CI footgun. Use NetworkManager/networkd config.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1785",
		Title:    "Error on `ufw default allow` — flips host firewall from deny-by-default to allow-by-default",
		Severity: SeverityError,
		Description: "`ufw default allow incoming` (or `allow outgoing`, `allow routed`) changes " +
			"the chain's baseline verdict — instead of only what you explicitly opened, every " +
			"port that does not have a matching `deny` rule is accepted. On an internet-facing " +
			"host this is effectively \"turn the firewall off\", and the effect survives reboots " +
			"because the default is persisted to `/etc/default/ufw`. Restore with `ufw default " +
			"deny incoming` and add narrow `ufw allow <port>` rules for the services that " +
			"actually need to be reachable.",
		Check: checkZC1785,
	})
}

func checkZC1785(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ufw" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	if cmd.Arguments[0].String() != "default" {
		return nil
	}

	verdict := cmd.Arguments[1].String()
	if verdict != "allow" {
		return nil
	}

	direction := "incoming"
	if len(cmd.Arguments) >= 3 {
		d := cmd.Arguments[2].String()
		if d == "incoming" || d == "outgoing" || d == "routed" {
			direction = d
		}
	}

	return []Violation{{
		KataID: "ZC1785",
		Message: "`ufw default allow " + direction + "` flips the firewall baseline to " +
			"accept every port that is not explicitly denied. Restore with `ufw default " +
			"deny incoming` and add narrow `ufw allow <port>` rules for the services that " +
			"must be reachable.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

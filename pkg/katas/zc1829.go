// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1829",
		Title:    "Warn on `tailscale down` / `wg-quick down` / `nmcli con down` — drops the VPN that may carry the SSH session",
		Severity: SeverityWarning,
		Description: "A script that closes the VPN tunnel from within a remote session cuts " +
			"itself off whenever the admin SSH rides over that tunnel. `tailscale down`, " +
			"`wg-quick down WG0`, `openvpn` teardown, and `nmcli connection down NAME` all " +
			"tear the link down in place with no grace or rollback. Schedule the teardown " +
			"behind `systemd-run --on-active=30s --unit=recover <cmd to bring it back up>` " +
			"so the VPN is back before the unit expires, or run the command from the " +
			"host's console / out-of-band path rather than over the VPN itself.",
		Check: checkZC1829,
	})
}

func checkZC1829(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "tailscale":
		if len(cmd.Arguments) > 0 && cmd.Arguments[0].String() == "down" {
			return zc1829Hit(cmd, "tailscale down")
		}
	case "wg-quick":
		if len(cmd.Arguments) > 0 && cmd.Arguments[0].String() == "down" {
			return zc1829Hit(cmd, "wg-quick down")
		}
	case "nmcli":
		// `nmcli connection down <name>` / `nmcli con down <name>`.
		if len(cmd.Arguments) >= 2 {
			first := cmd.Arguments[0].String()
			if (first == "connection" || first == "con" || first == "c") &&
				cmd.Arguments[1].String() == "down" {
				return zc1829Hit(cmd, "nmcli connection down")
			}
		}
	}
	return nil
}

func zc1829Hit(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1829",
		Message: "`" + what + "` tears down the VPN — if the SSH session rides on it, " +
			"the script cuts itself off with no rollback. Schedule recovery via " +
			"`systemd-run --on-active=30s`, or run from console / out-of-band.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

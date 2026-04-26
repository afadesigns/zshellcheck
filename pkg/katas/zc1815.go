// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1815NetUnits = map[string]bool{
	"NetworkManager":           true,
	"NetworkManager.service":   true,
	"systemd-networkd":         true,
	"systemd-networkd.service": true,
	"networking":               true,
	"networking.service":       true,
	"network":                  true,
	"network.service":          true,
	"wpa_supplicant":           true,
	"wpa_supplicant.service":   true,
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1815",
		Title:    "Warn on `systemctl restart NetworkManager` / `systemd-networkd` — drops the SSH session",
		Severity: SeverityWarning,
		Description: "Restarting the network manager from an SSH session tears down every active " +
			"connection the daemon supervises, including the one the script is running over. " +
			"The script freezes, the client sees a broken pipe, and recovery usually requires " +
			"console access. Route the change through `nmcli connection reload` + `nmcli " +
			"connection up <name>` (NetworkManager), `networkctl reload` (systemd-networkd), " +
			"or schedule the restart behind `systemd-run --on-active=30s` with a rollback " +
			"timer that re-enables the previous config if SSH does not reconnect.",
		Check: checkZC1815,
	})
}

func checkZC1815(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "systemctl" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	action := cmd.Arguments[0].String()
	if action != "restart" && action != "stop" && action != "reload-or-restart" {
		return nil
	}
	for _, arg := range cmd.Arguments[1:] {
		unit := strings.Trim(arg.String(), "\"'")
		if zc1815NetUnits[unit] {
			return []Violation{{
				KataID: "ZC1815",
				Message: "`systemctl " + action + " " + unit + "` drops every " +
					"connection the manager supervises — the SSH session freezes. " +
					"Use `nmcli connection reload` / `networkctl reload`, or a " +
					"`systemd-run --on-active=30s` rollback.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

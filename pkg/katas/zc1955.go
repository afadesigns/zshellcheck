package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1955",
		Title:    "Warn on `rfkill block all` / `block wifi|bluetooth|wwan` — disables every radio, cuts wireless",
		Severity: SeverityWarning,
		Description: "`rfkill block all` toggles the soft-kill switch on every radio the kernel " +
			"registered — WiFi, Bluetooth, WWAN, NFC, GPS, UWB — so the host drops off the " +
			"network in one call. A follow-up `rfkill unblock all` takes seconds to a minute " +
			"on some drivers and requires the operator to be physically present or have a " +
			"cellular fallback. Scope the block to a specific type (e.g. `rfkill block " +
			"bluetooth`) and schedule via `at now + 5 minutes ... rfkill unblock all` so the " +
			"host recovers on its own.",
		Check: checkZC1955,
	})
}

func checkZC1955(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "rfkill" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	if cmd.Arguments[0].String() != "block" {
		return nil
	}
	target := cmd.Arguments[1].String()
	if target != "all" && target != "wifi" && target != "wlan" &&
		target != "bluetooth" && target != "wwan" {
		return nil
	}
	return []Violation{{
		KataID: "ZC1955",
		Message: "`rfkill block " + target + "` hard-downs the radio(s) — host drops " +
			"off the network in one call. Scope to the radio type that really needs it " +
			"and schedule an `at now + N minutes` unblock for self-recovery.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

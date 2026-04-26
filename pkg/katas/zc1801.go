// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1801",
		Title:    "Warn on `fwupdmgr update` / `install` — mid-flash interruption can brick firmware",
		Severity: SeverityWarning,
		Description: "`fwupdmgr update`, `fwupdmgr upgrade`, and `fwupdmgr install FIRMWARE` push " +
			"new firmware into BIOS / UEFI, SSD, Thunderbolt controller, NIC, or dock " +
			"microcontroller. Most of those devices have no A/B rollback — an interrupted " +
			"flash (power cut, unexpected reboot, PSU toggle) leaves the chip in an " +
			"unbootable state that needs vendor-recovery hardware. Run from a battery-backed " +
			"session, mask reboot triggers with `systemd-inhibit`, pin the power supply, and " +
			"verify the update history with `fwupdmgr get-history` once the device returns.",
		Check: checkZC1801,
	})
}

func checkZC1801(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "fwupdmgr" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	switch cmd.Arguments[0].String() {
	case "update", "upgrade", "install", "reinstall", "downgrade":
		return []Violation{{
			KataID: "ZC1801",
			Message: "`fwupdmgr " + cmd.Arguments[0].String() + "` flashes firmware — a " +
				"mid-write interruption can brick BIOS, SSD, Thunderbolt, or NIC " +
				"microcontrollers. Inhibit reboot triggers (`systemd-inhibit`) and " +
				"ensure battery / UPS before running.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}
	return nil
}

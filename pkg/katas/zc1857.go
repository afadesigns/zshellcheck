// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1857",
		Title:    "Error on `cloud-init clean` — wipes boot state, next reboot re-provisions the host",
		Severity: SeverityError,
		Description: "`cloud-init clean` (and variants `--logs`, `--reboot`, `--machine-id`) " +
			"removes every marker under `/var/lib/cloud/` and `/var/log/cloud-init*`, " +
			"which tells cloud-init to re-run from scratch on the next boot. That run " +
			"re-imports the image-builder's user-data: regenerates SSH host keys, resets " +
			"the hostname, replaces `/etc/fstab` entries the operator may have edited, " +
			"and (with `--reboot`) triggers the replay immediately. In a maintenance " +
			"script this silently erases everything the operator configured after " +
			"first-boot. Keep the command out of automation; if you truly need to " +
			"re-seed an instance, snapshot state first and run the command interactively.",
		Check: checkZC1857,
	})
}

func checkZC1857(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cloud-init" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	if cmd.Arguments[0].String() != "clean" {
		return nil
	}
	return []Violation{{
		KataID: "ZC1857",
		Message: "`cloud-init clean` wipes `/var/lib/cloud/` boot state — the next " +
			"reboot re-runs the user-data and overwrites operator changes " +
			"(SSH host keys, hostname, `/etc/fstab`). Run interactively only.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

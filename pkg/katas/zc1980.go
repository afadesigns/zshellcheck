// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1980",
		Title:    "Error on `udevadm trigger --action=remove` — replays `remove` uevents, detaches live devices",
		Severity: SeverityError,
		Description: "`udevadm trigger --action=remove` (also spelled `-c remove`) walks " +
			"`/sys` and synthesises a `remove` uevent for every matching device. The " +
			"kernel reacts as if every matched disk, NIC, GPU, or USB node was " +
			"physically yanked — SATA controllers detach drives that back mounted " +
			"filesystems, netdevs disappear mid-session, and `systemd-udevd` fires " +
			"per-device cleanup rules it was never meant to run on a live host. The " +
			"normal way to replay `add`/`change` events after a rules edit is " +
			"`udevadm control --reload` followed by `udevadm trigger` with the default " +
			"action (`change`); scope any `--action=remove` to a specific device " +
			"subsystem with `--subsystem-match=` + `--attr-match=` and test on a " +
			"non-production box first.",
		Check: checkZC1980,
	})
}

func checkZC1980(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "udevadm" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "trigger" {
		return nil
	}
	for i, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if strings.HasPrefix(v, "--action=") {
			if strings.TrimPrefix(v, "--action=") == "remove" {
				return zc1980Hit(cmd)
			}
		}
		if v == "-c" || v == "--action" {
			if i+2 < len(cmd.Arguments) {
				next := cmd.Arguments[i+2].String()
				if next == "remove" {
					return zc1980Hit(cmd)
				}
			}
		}
	}
	return nil
}

func zc1980Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1980",
		Message: "`udevadm trigger --action=remove` replays `remove` uevents across " +
			"`/sys` — SATA/NIC/GPU nodes detach on a live host. Reload rules " +
			"with `udevadm control --reload`; scope with `--subsystem-match=`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

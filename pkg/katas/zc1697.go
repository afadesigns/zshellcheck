// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1697",
		Title:    "Info: `cryptsetup open --allow-discards` — TRIM pass-through leaks free-sector map",
		Severity: SeverityInfo,
		Description: "`--allow-discards` tells dm-crypt to forward TRIM/DISCARD commands from " +
			"the filesystem to the underlying SSD. The performance and wear-levelling gains " +
			"are real, but so is the side effect: an attacker with raw-device access can " +
			"read the free-sector map and see which blocks are empty — enough to fingerprint " +
			"partition layouts, distinguish encrypted-full-volume from encrypted-sparse-" +
			"content cases, and defeat plausible-deniability scenarios. If the threat model " +
			"includes offline-disk inspection, drop `--allow-discards` and accept the perf " +
			"hit; otherwise keep the flag but state the trade-off in the runbook.",
		Check: checkZC1697,
	})
}

func checkZC1697(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "cryptsetup" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "--allow-discards" {
			return []Violation{{
				KataID: "ZC1697",
				Message: "`cryptsetup --allow-discards` leaks free-sector layout to anyone " +
					"with raw-device access — drop it if offline-disk inspection is in " +
					"scope, or document the trade-off.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityInfo,
			}}
		}
	}
	return nil
}

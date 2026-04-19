package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1792",
		Title:    "Warn on `btrfs subvolume delete` / `btrfs device remove` — unrecoverable btrfs data loss",
		Severity: SeverityWarning,
		Description: "`btrfs subvolume delete PATH` unlinks the subvolume and drops all of its " +
			"extents once cleanup completes — on Snapper / Timeshift systems the argument is " +
			"often a snapshot that is the only remaining copy of pre-incident state. " +
			"`btrfs device remove DEV POOL` moves the stored chunks off DEV before detaching " +
			"it; wrong device, mid-rebalance failure, or insufficient free space across the " +
			"remaining members puts the filesystem into degraded mode with no automatic " +
			"rollback. Keep a fresh `btrfs subvolume list`/`btrfs device usage` snapshot and " +
			"confirm the target explicitly before running either command in automation.",
		Check: checkZC1792,
	})
}

func checkZC1792(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "btrfs" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}

	sub0 := cmd.Arguments[0].String()
	sub1 := cmd.Arguments[1].String()

	switch {
	case sub0 == "subvolume" && sub1 == "delete":
		return zc1792Hit(cmd, "btrfs subvolume delete")
	case sub0 == "device" && sub1 == "remove":
		return zc1792Hit(cmd, "btrfs device remove")
	}
	return nil
}

func zc1792Hit(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1792",
		Message: "`" + what + "` drops btrfs state with no automatic rollback — " +
			"snapshots vanish on `subvolume delete`, and `device remove` can leave " +
			"the filesystem degraded. Confirm the target explicitly.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1646",
		Title:    "Warn on `btrfs check --repair` / `xfs_repair -L` — last-resort recovery, may worsen damage",
		Severity: SeverityWarning,
		Description: "Both commands are destructive last-resort recovery. `btrfs check " +
			"--repair` explicitly warns in its man page that it \"may cause additional " +
			"filesystem damage\" and the btrfs developers ask users to try `btrfs scrub` and " +
			"read-only `btrfs check` first. `xfs_repair -L` zeroes the log, dropping any " +
			"uncommitted transactions and the data they held. In both cases snapshot the " +
			"underlying block device before running, so the attempt is reversible.",
		Check: checkZC1646,
	})
}

func checkZC1646(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value == "btrfs" {
		if len(cmd.Arguments) < 2 || cmd.Arguments[0].String() != "check" {
			return nil
		}
		for _, arg := range cmd.Arguments[1:] {
			if arg.String() == "--repair" {
				return []Violation{{
					KataID: "ZC1646",
					Message: "`btrfs check --repair` may worsen damage — try `btrfs scrub` " +
						"and read-only `btrfs check` first, and snapshot the block device " +
						"before running.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
		return nil
	}
	if ident.Value == "xfs_repair" {
		for _, arg := range cmd.Arguments {
			if arg.String() == "-L" {
				return []Violation{{
					KataID: "ZC1646",
					Message: "`xfs_repair -L` zeroes the log — uncommitted transactions are " +
						"lost. Snapshot the block device first; mount read-only and read " +
						"the log if possible.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityWarning,
				}}
			}
		}
	}
	return nil
}

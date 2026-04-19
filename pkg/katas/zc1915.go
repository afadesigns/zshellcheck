package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1915",
		Title:    "Error on `mdadm --zero-superblock` / `--stop` — drops RAID metadata or live array",
		Severity: SeverityError,
		Description: "`mdadm --zero-superblock $DEV` wipes the MD superblock from a member — the " +
			"array forgets the device exists and a subsequent `--create` with the wrong layout " +
			"permanently scrambles the data. `mdadm --stop $MD` (or `-S`) halts a live array " +
			"from underneath whatever is mounted on it; if root or `/boot` lives there the host " +
			"panics on the next fsync. Run `mdadm --examine` first, snapshot the superblock " +
			"with `mdadm --detail --export`, and keep both calls behind a runbook rather than " +
			"an automated script.",
		Check: checkZC1915,
	})
}

func checkZC1915(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	// Parser caveat: `mdadm --zero-superblock $DEV` mangles the command
	// name to `zero-superblock`.
	switch ident.Value {
	case "zero-superblock":
		return zc1915Hit(cmd, "mdadm --zero-superblock")
	case "mdadm":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			switch v {
			case "--zero-superblock":
				return zc1915Hit(cmd, "mdadm --zero-superblock")
			case "-S", "--stop":
				return zc1915Hit(cmd, "mdadm "+v)
			case "--remove":
				return zc1915Hit(cmd, "mdadm --remove")
			}
		}
	}
	return nil
}

func zc1915Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1915",
		Message: "`" + form + "` drops RAID metadata or halts a live array — mounted root " +
			"or /boot panics the host; a stale superblock scrambles data on next `--create`. " +
			"Snapshot `mdadm --detail --export` first and keep behind a runbook.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

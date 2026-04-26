// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1663",
		Title:    "Warn on `tune2fs -c 0` / `-i 0` — disables periodic filesystem checks",
		Severity: SeverityWarning,
		Description: "`tune2fs -c 0` (mount count) and `tune2fs -i 0` (time interval) disable " +
			"the ext2/3/4 periodic-fsck machinery so the filesystem only gets checked after a " +
			"dirty unmount or a manual `fsck -f`. For desktops the nag is annoying; for " +
			"long-lived servers it is the last line of defence against silent metadata " +
			"corruption. Lower the cadence if the default is too aggressive (`tune2fs -c 30`, " +
			"`-i 3m`) rather than turning it off, and schedule an offline `fsck` on a cadence " +
			"you can defend.",
		Check: checkZC1663,
	})
}

func checkZC1663(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "tune2fs" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v != "-c" && v != "-i" {
			continue
		}
		if i+1 >= len(cmd.Arguments) {
			continue
		}
		next := cmd.Arguments[i+1].String()
		if next != "0" {
			continue
		}
		return []Violation{{
			KataID: "ZC1663",
			Message: "`tune2fs " + v + " 0` disables periodic fsck on the filesystem — " +
				"lower the cadence (e.g. `-c 30` / `-i 3m`) instead of turning it off.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}
	return nil
}

// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1981",
		Title:    "Warn on `exec -a NAME cmd` — replaces `argv[0]`, hides the real binary from `ps`",
		Severity: SeverityWarning,
		Description: "`exec -a NAME $BIN` tells Zsh to set `argv[0]` of the `exec`'d process " +
			"to `NAME` instead of the actual program path. `ps`, `top`, `proc`-based " +
			"audit tools, and systemd's unit accounting all see `NAME` — the real " +
			"binary on disk is only discoverable from `/proc/PID/exe`, which most " +
			"monitoring does not read. The feature has legitimate uses (login shells " +
			"spelling themselves `-zsh` so tty/shell detection works) but also makes a " +
			"great disguise for a reverse shell or a cron-triggered helper. Keep " +
			"`exec -a` out of production scripts unless the intent is documented; " +
			"prefer running the binary at its real path so operators can match process " +
			"name to on-disk file.",
		Check: checkZC1981,
	})
}

func checkZC1981(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "exec" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		if arg.String() == "-a" {
			return []Violation{{
				KataID: "ZC1981",
				Message: "`exec -a NAME` sets `argv[0]` to `NAME` — `ps`/`top`/audit " +
					"rules see the alias, not the real binary. Keep out of production " +
					"scripts unless the alias is documented.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1493",
		Title:    "Warn on `watch -n 0` — zero-interval watch spins CPU",
		Severity: SeverityWarning,
		Description: "`watch -n 0` (or `-n 0.0` / `-n .0`) tells `watch` to re-run the command " +
			"with no delay, which immediately pins a core to 100% and usually saturates the " +
			"terminal emulator too. Pick a realistic interval (`-n 1`, `-n 2`, `-n 0.5`) — or " +
			"if you truly want tight polling, use a dedicated event API (`inotifywait`, " +
			"`systemd.path` unit, `journalctl -f`).",
		Check: checkZC1493,
	})
}

func checkZC1493(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "watch" {
		return nil
	}

	var prevN bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if prevN {
			prevN = false
			if v == "0" || v == "0.0" || v == "0.00" || v == ".0" || v == ".00" {
				return zc1493Violation(cmd, v)
			}
		}
		if v == "-n" || v == "--interval" {
			prevN = true
			continue
		}
		if v == "-n0" || v == "-n0.0" || v == "--interval=0" || v == "--interval=0.0" {
			return zc1493Violation(cmd, v)
		}
	}
	return nil
}

func zc1493Violation(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1493",
		Message: "`watch -n " + what + "` pins a core at 100% and saturates the terminal. " +
			"Use a realistic interval or an event-driven API (inotifywait / journalctl -f).",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

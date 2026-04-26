// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1939",
		Title:    "Error on `reboot -f` / `halt -f` / `poweroff -f` — skips shutdown sequence, no graceful service stop",
		Severity: SeverityError,
		Description: "`reboot -f`, `halt -f`, and `poweroff -f` short-circuit the systemd " +
			"shutdown graph — no `ExecStop=`, no `DefaultDependencies=`, no filesystem sync, " +
			"no Before/After ordering. The kernel's `reboot(2)` fires immediately and every " +
			"dirty buffer that was not yet flushed is lost. Journal writes stop mid-line, " +
			"databases on the host replay from the last checkpoint, and anything that needed a " +
			"clean unmount (LUKS, NFS, cephfs) logs a dirty state. Use plain `systemctl " +
			"reboot` / `shutdown -r +N`, and reserve `-f` for recovery when the normal path is " +
			"already wedged.",
		Check: checkZC1939,
	})
}

var zc1939ForceFlags = map[string]bool{
	"-f":      true,
	"--force": true,
}

func checkZC1939(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "reboot" && ident.Value != "halt" && ident.Value != "poweroff" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-f" || v == "--force" {
			return zc1939Hit(cmd, ident.Value+" "+v)
		}
		if len(v) >= 2 && v[0] == '-' && v[1] != '-' {
			for i := 1; i < len(v); i++ {
				if v[i] == 'f' {
					return zc1939Hit(cmd, ident.Value+" -f")
				}
			}
		}
	}
	return nil
}

func zc1939Hit(cmd *ast.SimpleCommand, form string) []Violation {
	line, col := FlagArgPosition(cmd, zc1939ForceFlags)
	return []Violation{{
		KataID: "ZC1939",
		Message: "`" + form + "` fires `reboot(2)` immediately — no `ExecStop=`, no " +
			"filesystem sync, no clean unmount. Databases replay from last checkpoint. " +
			"Use `systemctl reboot` / `shutdown -r +N`; reserve `-f` for wedged recovery.",
		Line:   line,
		Column: col,
		Level:  SeverityError,
	}}
}

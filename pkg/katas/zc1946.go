// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1946",
		Title:    "Warn on `unsetopt HUP` — background jobs keep running after shell exit",
		Severity: SeverityWarning,
		Description: "Zsh's `HUP` option (on by default) sends `SIGHUP` to each running child job " +
			"when the shell exits, letting them wind down cleanly. `unsetopt HUP` / " +
			"`setopt NO_HUP` disables that, so long pipelines, `sleep` loops, and " +
			"user-spawned daemons live on — `ps aux` accumulates orphaned workers across " +
			"logouts and resource consumption creeps up. If a specific job really needs to " +
			"outlive the shell, use `disown` or `systemd-run --scope` on that one invocation; " +
			"leave `HUP` on globally.",
		Check: checkZC1946,
	})
}

func checkZC1946(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var disabling bool
	switch ident.Value {
	case "unsetopt":
		disabling = true
	case "setopt":
		disabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1946Canonical(arg.String())
		switch v {
		case "HUP":
			if disabling {
				return zc1946Hit(cmd, "unsetopt HUP")
			}
		case "NOHUP":
			if !disabling {
				return zc1946Hit(cmd, "setopt NO_HUP")
			}
		}
	}
	return nil
}

func zc1946Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc1946Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1946",
		Message: "`" + form + "` stops the shell from `SIGHUP`-ing background jobs on " +
			"exit — long pipelines and spawned daemons outlive the session, orphans " +
			"accumulate. Use `disown` or `systemd-run --scope` on specific commands instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

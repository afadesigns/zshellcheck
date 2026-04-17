package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1597",
		Title:    "Warn on `systemd-run -p User=root` — launches arbitrary command with root privileges",
		Severity: SeverityWarning,
		Description: "`systemd-run` submits a transient unit to systemd. With `-p User=root` " +
			"(or `User=0`) the unit runs as root — bypassing the usual `sudo` audit path in " +
			"`/var/log/auth.log`. On hosts where the caller's polkit / dbus rules allow the " +
			"operation, this is effectively privilege escalation by a different name. Prefer " +
			"explicit `sudo` so the invocation is logged, or pre-provision a dedicated systemd " +
			"unit that names the exact command it can run.",
		Check: checkZC1597,
	})
}

func checkZC1597(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "systemd-run" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "User=root" || v == "User=0" {
			return []Violation{{
				KataID: "ZC1597",
				Message: "`systemd-run -p " + v + "` runs arbitrary commands as root via " +
					"systemd — bypasses the `sudo` audit path. Prefer explicit `sudo` or a " +
					"fixed systemd unit.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

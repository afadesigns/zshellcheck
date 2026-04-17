package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1512",
		Title:    "Style: `service <unit> <verb>` — use `systemctl <verb> <unit>` on systemd hosts",
		Severity: SeverityStyle,
		Description: "`service` is the SysV init compatibility wrapper. On a systemd-managed " +
			"host (every mainstream distro since ~2016) it translates to `systemctl` anyway, " +
			"but reverses argument order, loses `--user` scope, ignores unit templating, and " +
			"can't restart sockets or timers. Prefer `systemctl start|stop|restart|reload " +
			"<unit>` for consistency across scripts and interactive shells.",
		Check: checkZC1512,
	})
}

func checkZC1512(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "service" {
		return nil
	}

	// Needs at least <unit> <verb>.
	if len(cmd.Arguments) < 2 {
		return nil
	}
	verb := cmd.Arguments[1].String()
	switch verb {
	case "start", "stop", "restart", "reload", "status", "force-reload", "try-restart":
	default:
		return nil
	}

	unit := cmd.Arguments[0].String()
	return []Violation{{
		KataID: "ZC1512",
		Message: "`service " + unit + " " + verb + "` — prefer `systemctl " + verb + " " +
			unit + "` for consistency with other systemd commands.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

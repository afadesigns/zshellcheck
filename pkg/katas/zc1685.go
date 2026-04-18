package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1685",
		Title:    "Info: `sleep infinity` — container keep-alive pattern that ignores SIGTERM",
		Severity: SeverityInfo,
		Description: "`sleep infinity` is most often used as a container or systemd-unit keep-" +
			"alive. Problem: GNU `sleep` does not install a SIGTERM handler, so when " +
			"`docker stop` / `systemctl stop` sends SIGTERM the process sits unresponsive " +
			"until the grace period expires and SIGKILL lands. The orchestrator reports a " +
			"hung stop, logs look wrong, and any cleanup registered on signal handlers in " +
			"a wrapping shell never runs. Replace with `exec tail -f /dev/null` (signal-" +
			"handles cleanly) or front with `tini` / `dumb-init` when PID 1 must stay.",
		Check: checkZC1685,
	})
}

func checkZC1685(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "sleep" {
		return nil
	}
	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "infinity" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1685",
		Message: "`sleep infinity` does not trap SIGTERM — the orchestrator hangs until " +
			"SIGKILL. Use `exec tail -f /dev/null` or front with `tini` / `dumb-init`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}

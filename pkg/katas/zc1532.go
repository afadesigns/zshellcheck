package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1532",
		Title:    "Warn on `screen -dm` / `tmux new-session -d` — detached long-running session",
		Severity: SeverityWarning,
		Description: "Starting a detached screen/tmux session from a script puts a long-running " +
			"process outside the systemd supervisory tree: no logs in the journal, no cgroup " +
			"accounting, no restart-on-failure, no OOM scoring. It is also a common post- " +
			"compromise persistence technique because the session survives the initial shell " +
			"exit and hides in `ps -ef` as a short tmux/screen helper. For real long-running " +
			"work, write a systemd unit (user or system) and start it with `systemctl " +
			"[--user] start`.",
		Check: checkZC1532,
	})
}

func checkZC1532(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value == "screen" {
		for _, arg := range cmd.Arguments {
			v := arg.String()
			if v == "-dm" || v == "-dmS" {
				return zc1532Violation(cmd, "screen "+v)
			}
		}
	}
	if ident.Value == "tmux" && len(cmd.Arguments) >= 2 && cmd.Arguments[0].String() == "new-session" {
		for _, arg := range cmd.Arguments[1:] {
			if arg.String() == "-d" {
				return zc1532Violation(cmd, "tmux new-session -d")
			}
		}
	}
	return nil
}

func zc1532Violation(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1532",
		Message: "`" + what + "` backgrounds work outside systemd — no journal, no cgroup, " +
			"common persistence technique. Use a systemd unit instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

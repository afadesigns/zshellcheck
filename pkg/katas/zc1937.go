package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1937",
		Title:    "Warn on `tmux kill-server` / `tmux kill-session` — tears down every detached process inside",
		Severity: SeverityWarning,
		Description: "`tmux kill-server` terminates the whole tmux daemon, `tmux kill-session -t NAME` " +
			"drops one named session, and `screen -X quit` does the screen equivalent. Anything " +
			"the operator parked inside — a long-running build, a `tail -F` on production " +
			"logs, a held `sudo` token, a port-forward — dies with the session, and the " +
			"detached processes get `SIGHUP`'d with no cleanup. Use `tmux kill-window -t …` for " +
			"surgical removal, send `SIGTERM` to the specific backend, or rely on `systemd-run " +
			"--scope` for workloads that should survive terminal churn.",
		Check: checkZC1937,
	})
}

func checkZC1937(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value == "tmux" && len(cmd.Arguments) >= 1 {
		sub := cmd.Arguments[0].String()
		if sub == "kill-server" || sub == "kill-session" {
			return zc1937Hit(cmd, "tmux "+sub)
		}
	}
	if ident.Value == "screen" {
		for i, arg := range cmd.Arguments {
			if arg.String() == "-X" && i+1 < len(cmd.Arguments) &&
				cmd.Arguments[i+1].String() == "quit" {
				return zc1937Hit(cmd, "screen -X quit")
			}
		}
	}
	return nil
}

func zc1937Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1937",
		Message: "`" + form + "` tears down every detached process inside the session — " +
			"builds, log tails, port-forwards get `SIGHUP`'d with no cleanup. Use " +
			"`kill-window` for surgical removal or `systemd-run --scope` for workloads.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1435",
		Title:    "Avoid `killall -9` / `killall -KILL` — force-kill by process name",
		Severity: SeverityWarning,
		Description: "`killall -9 name` sends SIGKILL to every process matching `name` — in " +
			"multi-user or containerized environments, this can hit unrelated processes that " +
			"happen to share the name. Prefer `killall -TERM` first (graceful), or kill by PID " +
			"after locating with `pgrep` / `pidof`.",
		Check: checkZC1435,
	})
}

func checkZC1435(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "killall" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-9" || v == "-KILL" || v == "-s" {
			return []Violation{{
				KataID: "ZC1435",
				Message: "`killall -9 name` force-kills every matching process, including " +
					"unrelated instances on multi-user or containerized hosts. Start with -TERM, " +
					"or kill by PID after `pgrep`/`pidof`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}

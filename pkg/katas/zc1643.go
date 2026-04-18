package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1643",
		Title:    "Style: `$(cat file)` — use `$(<file)` to skip the fork / exec",
		Severity: SeverityStyle,
		Description: "`$(cat FILE)` forks, execs `/usr/bin/cat`, reads FILE, writes the bytes " +
			"to the pipe, waits for the child. `$(<FILE)` is a shell builtin — it reads FILE " +
			"directly into the command-substitution buffer with no fork and no exec. In a hot " +
			"path the speedup is dramatic, and even in cold paths it avoids one of the most " +
			"common useless-use-of-cat patterns in review feedback.",
		Check: checkZC1643,
	})
}

func checkZC1643(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.Contains(v, "$(cat ") {
			return []Violation{{
				KataID: "ZC1643",
				Message: "`$(cat FILE)` forks cat just to read a file — use `$(<FILE)` " +
					"(shell builtin, no fork).",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}
	return nil
}

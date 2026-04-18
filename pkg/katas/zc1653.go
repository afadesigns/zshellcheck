package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1653",
		Title:    "Avoid `$BASHPID` — Bash-only; Zsh uses `$sysparams[pid]` from `zsh/system`",
		Severity: SeverityWarning,
		Description: "`$BASHPID` returns the PID of the current subshell (while `$$` returns " +
			"the parent shell's PID). In Zsh this parameter is not set — scripts that rely " +
			"on `$BASHPID` silently get an empty string and misbehave. After `zmodload " +
			"zsh/system`, Zsh exposes the current process PID as `$sysparams[pid]`, which " +
			"updates inside subshells just like Bash's `$BASHPID`.",
		Check: checkZC1653,
	})
}

func checkZC1653(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.Contains(v, "$BASHPID") || strings.Contains(v, "${BASHPID}") {
			return []Violation{{
				KataID: "ZC1653",
				Message: "`$BASHPID` is Bash-only. Use `$sysparams[pid]` after " +
					"`zmodload zsh/system`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

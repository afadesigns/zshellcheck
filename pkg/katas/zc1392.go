package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1392",
		Title:    "Avoid `$CHILD_MAX` — Bash-only; Zsh uses `limit` / `ulimit -u`",
		Severity: SeverityInfo,
		Description: "Bash's `$CHILD_MAX` reports the maximum number of exited child processes " +
			"Bash remembers. Zsh does not export this var. For current process limits use " +
			"`limit -s maxproc` or `ulimit -u` — but the exact Bash semantic is not mirrored.",
		Check: checkZC1392,
	})
}

func checkZC1392(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "echo" && ident.Value != "print" && ident.Value != "printf" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.Contains(v, "CHILD_MAX") {
			return []Violation{{
				KataID: "ZC1392",
				Message: "`$CHILD_MAX` is Bash-only. Zsh uses `limit -s maxproc` or `ulimit -u` " +
					"for process-count limits.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityInfo,
			}}
		}
	}

	return nil
}

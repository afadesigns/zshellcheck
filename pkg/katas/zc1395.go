package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1395",
		Title:    "Avoid `wait -n` — Bash 4.3+ only; Zsh `wait` on job IDs",
		Severity: SeverityWarning,
		Description: "Bash 4.3+ added `wait -n` (wait for any job to finish). Zsh's `wait` does " +
			"not accept `-n`; instead wait explicitly on job IDs or PIDs, or use `wait` with no " +
			"args (waits for all). For any-of semantics use `wait $pid1 $pid2; ...` in a loop.",
		Check: checkZC1395,
	})
}

func checkZC1395(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "wait" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-n" {
			return []Violation{{
				KataID: "ZC1395",
				Message: "`wait -n` is Bash 4.3+. Zsh's `wait` waits on specific PIDs/jobs or " +
					"(bare `wait`) all jobs. For any-child semantics, loop over PIDs with " +
					"individual `wait $pid` calls.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}

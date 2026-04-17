package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1430",
		Title:    "Prefer Zsh `zsh/sched` module over `at now` / `batch` for in-shell scheduling",
		Severity: SeverityStyle,
		Description: "`at`/`batch` schedule commands via the atd daemon — requires daemon " +
			"running, leaves a spool-file audit trail, and runs in a fresh environment. For " +
			"in-shell scheduling the Zsh `zsh/sched` module (`sched +1:00 cmd`) runs the " +
			"command from the current shell without the daemon dependency.",
		Check: checkZC1430,
	})
}

func checkZC1430(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "at" && ident.Value != "batch" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1430",
		Message: "Prefer Zsh `sched` (from `zsh/sched`) for in-shell scheduling instead of " +
			"`at`/`batch`. No daemon dependency, runs in the current shell's environment.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1167",
		Title:    "Avoid `timeout` command — use Zsh `TMOUT` or `zsh/sched`",
		Severity: SeverityStyle,
		Description: "`timeout` is not available on all systems (macOS lacks it by default). " +
			"Use Zsh `TMOUT` variable or `zmodload zsh/sched` for timeout functionality.",
		Check: checkZC1167,
	})
}

func checkZC1167(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "timeout" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1167",
		Message: "Avoid `timeout` — it's unavailable on macOS. Use Zsh `TMOUT` variable " +
			"or `zmodload zsh/sched` for portable timeout functionality.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

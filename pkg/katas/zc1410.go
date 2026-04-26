// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1410",
		Title:    "Avoid `compopt` — Bash programmable-completion modifier, not in Zsh",
		Severity: SeverityError,
		Description: "`compopt` tweaks Bash programmable-completion options for the current " +
			"completion. Zsh's compsys does not implement `compopt`; completion options are set " +
			"via `zstyle` / completion-function context instead.",
		Check: checkZC1410,
	})
}

func checkZC1410(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "compopt" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1410",
		Message: "`compopt` is a Bash-only completion builtin. Zsh compsys uses `zstyle` " +
			"(e.g. `zstyle ':completion:*' menu select`) for equivalent tuning.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

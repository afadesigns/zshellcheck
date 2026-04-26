// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1296",
		Title:    "Avoid `shopt` in Zsh — use `setopt`/`unsetopt` instead",
		Severity: SeverityWarning,
		Description: "`shopt` is a Bash builtin that does not exist in Zsh. Use `setopt` " +
			"or `unsetopt` to control Zsh shell options. Common Bash `shopt` options " +
			"have Zsh equivalents via `setopt`.",
		Check: checkZC1296,
	})
}

func checkZC1296(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "shopt" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1296",
		Message: "Avoid `shopt` in Zsh — it is a Bash builtin. Use `setopt`/`unsetopt` for Zsh shell options.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityWarning,
	}}
}

// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1337",
		Title:    "Avoid `fold` command — use Zsh `print -l` with `$COLUMNS`",
		Severity: SeverityStyle,
		Description: "`fold` wraps text to a specified width. Zsh provides `$COLUMNS` for " +
			"terminal width and `print -l` for line-by-line output, reducing " +
			"dependency on external commands.",
		Check: checkZC1337,
	})
}

func checkZC1337(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "fold" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1337",
		Message: "Consider Zsh `$COLUMNS` and `print` for text wrapping instead of `fold`.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityStyle,
	}}
}

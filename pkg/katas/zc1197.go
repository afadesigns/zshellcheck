// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1197",
		Title:    "Avoid `more` in scripts — use `cat` or pager check",
		Severity: SeverityStyle,
		Description: "`more` requires an interactive terminal and will hang in scripts. " +
			"Use `cat` for output or check `$TERM` before invoking a pager.",
		Check: checkZC1197,
	})
}

func checkZC1197(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "more" {
		return nil
	}

	hasFlags := false
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] == '-' {
			hasFlags = true
		}
	}

	if !hasFlags && len(cmd.Arguments) > 0 {
		return []Violation{{
			KataID: "ZC1197",
			Message: "Avoid `more` in scripts — it requires an interactive terminal. " +
				"Use `cat` for output or check `[[ -t 1 ]]` before paging.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}

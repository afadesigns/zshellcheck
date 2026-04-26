// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1316",
		Title:    "Avoid `caller` builtin — use `$funcfiletrace` in Zsh",
		Severity: SeverityWarning,
		Description: "`caller` is a Bash builtin that returns the call stack context. " +
			"Zsh provides `$funcfiletrace`, `$funcstack`, and `$funcsourcetrace` " +
			"for inspecting the call stack.",
		Check: checkZC1316,
	})
}

func checkZC1316(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "caller" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1316",
		Message: "Avoid `caller` in Zsh — it is a Bash builtin. Use `$funcfiletrace` and `$funcstack` instead.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityWarning,
	}}
}

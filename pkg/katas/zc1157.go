// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1157",
		Title:    "Avoid `strings` command — use Zsh `${(ps:\\0:)var}`",
		Severity: SeverityStyle,
		Description: "The `strings` command extracts printable strings from binaries. " +
			"For simple filtering, Zsh parameter expansion with `(ps:\\0:)` can split on null bytes.",
		Check: checkZC1157,
	})
}

func checkZC1157(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "strings" {
		return nil
	}

	// Only flag simple strings without special flags
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] == '-' {
			return nil
		}
	}

	if len(cmd.Arguments) != 1 {
		return nil
	}

	return []Violation{{
		KataID: "ZC1157",
		Message: "Consider Zsh parameter expansion for string extraction from variables. " +
			"`strings` is typically needed only for binary file analysis.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

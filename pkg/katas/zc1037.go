// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/token"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1037",
		Title: "Use 'print -r --' for variable expansion",
		Description: "Using 'echo' to print strings containing variables can lead to unexpected behavior " +
			"if the variable contains special characters or flags. A safer, more reliable alternative " +
			"is 'print -r --'.",
		Severity: SeverityStyle,
		Check:    checkZC1037,
		// Reuse ZC1092's `echo` → `print -r --` rewrite. The detector
		// here is a stricter subset (only fires when echo prints a
		// variable expansion) but the rewrite is identical; the
		// conflict resolver dedupes overlapping edits.
		Fix: fixZC1092,
	})
}

func checkZC1037(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	if cmd.Name.TokenLiteral() != "echo" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if ident, ok := arg.(*ast.Identifier); ok && ident.Token.Type == token.VARIABLE {
			return []Violation{
				{
					KataID:  "ZC1037",
					Message: "Use 'print -r --' instead of 'echo' to reliably print variable expansions.",
					Line:    cmd.Token.Line,
					Column:  cmd.Token.Column,
					Level:   SeverityStyle,
				},
			}
		}
		if str, ok := arg.(*ast.StringLiteral); ok && strings.Contains(str.Value, "$") {
			return []Violation{
				{
					KataID:  "ZC1037",
					Message: "Use 'print -r --' instead of 'echo' to reliably print variable expansions.",
					Line:    cmd.Token.Line,
					Column:  cmd.Token.Column,
					Level:   SeverityStyle,
				},
			}
		}
	}

	return nil
}

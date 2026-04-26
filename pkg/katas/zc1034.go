// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.ExpressionStatementNode, Kata{
		ID:    "ZC1034",
		Title: "Use `command -v` instead of `which`",
		Description: "`which` is an external command and may not be available or consistent across all " +
			"systems. `command -v` is a POSIX standard and a shell builtin, making it more portable " +
			"and reliable for checking if a command exists.",
		Severity: SeverityStyle,
		Check:    checkZC1034,
		Fix:      fixZC1034,
	})
}

// fixZC1034 rewrites `which` to `command -v` at the command name
// position inside an ExpressionStatement. Single replacement — arguments
// stay untouched.
func fixZC1034(node ast.Node, v Violation, source []byte) []FixEdit {
	es, ok := node.(*ast.ExpressionStatement)
	if !ok {
		return nil
	}
	cmd, ok := es.Expression.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if cmd.Name == nil || cmd.Name.TokenLiteral() != "which" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("which"),
		Replace: "command -v",
	}}
}

func checkZC1034(node ast.Node) []Violation {
	violations := []Violation{}

	if es, ok := node.(*ast.ExpressionStatement); ok {
		if cmd, ok := es.Expression.(*ast.SimpleCommand); ok {
			if cmd.Name.TokenLiteral() == "which" {
				violations = append(violations, Violation{
					KataID:  "ZC1034",
					Message: "Use `command -v` instead of `which` for portability.",
					Line:    cmd.Token.Line,
					Column:  cmd.Token.Column,
					Level:   SeverityStyle,
				})
			}
		}
	}

	return violations
}

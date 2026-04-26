// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.CommandSubstitutionNode, Kata{
		ID:    "ZC1015",
		Title: "Use `$(...)` for command substitution instead of backticks",
		Description: "The `$(...)` syntax is the modern, recommended way to perform command substitution. " +
			"It is more readable and can be nested easily, unlike backticks.",
		Severity: SeverityStyle,
		Check:    checkZC1015,
		Fix:      fixZC1002,
	})
}

func checkZC1015(node ast.Node) []Violation {
	violations := []Violation{}

	if cs, ok := node.(*ast.CommandSubstitution); ok {
		violations = append(violations, Violation{
			KataID:  "ZC1015",
			Message: "Use `$(...)` for command substitution instead of backticks.",
			Line:    cs.Token.Line,
			Column:  cs.Token.Column,
			Level:   SeverityStyle,
		})
	}

	return violations
}

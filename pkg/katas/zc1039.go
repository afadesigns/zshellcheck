// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1039",
		Title: "Avoid `rm` with root path",
		Description: "Running `rm` on the root directory `/` is dangerous. " +
			"Ensure you are not deleting the entire filesystem.",
		Severity: SeverityWarning,
		Check:    checkZC1039,
	})
}

func checkZC1039(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	// Check if command is rm
	if cmdName, ok := cmd.Name.(*ast.Identifier); !ok || cmdName.Value != "rm" {
		return nil
	}

	violations := []Violation{}

	for _, arg := range cmd.Arguments {
		// Bare `/` argument arrives as an Identifier with Value "/"
		// after the SLASH prefix registration; quoted forms arrive
		// as StringLiteral. Cover both shapes.
		var val string
		var line, col int
		switch n := arg.(type) {
		case *ast.StringLiteral:
			val = strings.Trim(n.Value, "\"'")
			line, col = n.Token.Line, n.Token.Column
		case *ast.Identifier:
			val = n.Value
			line, col = n.Token.Line, n.Token.Column
		}
		if val == "/" {
			violations = append(violations, Violation{
				KataID:  "ZC1039",
				Message: "Avoid `rm` on the root directory `/`. This is highly dangerous.",
				Line:    line,
				Column:  col,
				Level:   SeverityWarning,
			})
		}
	}

	return violations
}

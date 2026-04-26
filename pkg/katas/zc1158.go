// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1158",
		Title:    "Avoid `chown -R` without `--no-dereference`",
		Severity: SeverityWarning,
		Description: "`chown -R` follows symlinks by default, potentially changing ownership " +
			"outside the intended directory. Use `--no-dereference` or `-h` to avoid this.",
		Check: checkZC1158,
	})
}

func checkZC1158(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "chown" {
		return nil
	}

	hasRecursive := false
	hasSafe := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-R" {
			hasRecursive = true
		}
		if val == "-h" || val == "--no-dereference" {
			hasSafe = true
		}
	}

	if hasRecursive && !hasSafe {
		return []Violation{{
			KataID: "ZC1158",
			Message: "Use `chown -Rh` or `chown -R --no-dereference` to prevent following " +
				"symlinks during recursive ownership changes.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}

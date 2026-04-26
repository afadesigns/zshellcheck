// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1280",
		Title:    "Use `Zsh ${var:e}` instead of shell expansion to extract file extension",
		Severity: SeverityStyle,
		Description: "Zsh provides the `:e` (extension) modifier for parameter expansion which " +
			"extracts the file extension, avoiding complex shell patterns or external tools.",
		Check: checkZC1280,
	})
}

func checkZC1280(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cut" {
		return nil
	}

	hasDot := false
	hasField := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-d." || val == "-d" {
			hasDot = true
		}
		if val == "-f2" {
			hasField = true
		}
		if val == "." {
			hasDot = true
		}
	}

	if hasDot && hasField {
		return []Violation{{
			KataID:  "ZC1280",
			Message: "Use Zsh parameter expansion `${var:e}` to extract the file extension instead of `cut -d. -f2`.",
			Line:    cmd.Token.Line,
			Column:  cmd.Token.Column,
			Level:   SeverityStyle,
		}}
	}

	return nil
}

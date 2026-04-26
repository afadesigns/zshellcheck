// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1064",
		Title: "Prefer `command -v` over `type`",
		Description: "`type` output format varies and is not POSIX standard for checking existence. " +
			"`command -v` is quieter and standard.",
		Severity: SeverityInfo,
		Check:    checkZC1064,
		Fix:      fixZC1064,
	})
}

// fixZC1064 rewrites `type` to `command -v` at the command name
// position. Single replacement — arguments stay untouched.
func fixZC1064(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "type" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("type"),
		Replace: "command -v",
	}}
}

func checkZC1064(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	if name, ok := cmd.Name.(*ast.Identifier); ok && name.Value == "type" {
		return []Violation{{
			KataID:  "ZC1064",
			Message: "Prefer `command -v` over `type`. `type` output is not stable/standard for checking command existence.",
			Line:    name.Token.Line,
			Column:  name.Token.Column,
			Level:   SeverityInfo,
		}}
	}

	return nil
}

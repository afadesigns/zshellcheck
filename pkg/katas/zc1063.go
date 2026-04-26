// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:          "ZC1063",
		Title:       "Prefer `grep -F` over `fgrep`",
		Description: "`fgrep` is deprecated. Use `grep -F` instead.",
		Severity:    SeverityInfo,
		Check:       checkZC1063,
		Fix:         fixZC1063,
	})
}

// fixZC1063 rewrites `fgrep` to `grep -F` at the command name
// position. Single replacement — arguments stay untouched.
func fixZC1063(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "fgrep" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("fgrep"),
		Replace: "grep -F",
	}}
}

func checkZC1063(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	if name, ok := cmd.Name.(*ast.Identifier); ok && name.Value == "fgrep" {
		return []Violation{{
			KataID:  "ZC1063",
			Message: "`fgrep` is deprecated. Use `grep -F` instead.",
			Line:    name.Token.Line,
			Column:  name.Token.Column,
			Level:   SeverityInfo,
		}}
	}

	return nil
}

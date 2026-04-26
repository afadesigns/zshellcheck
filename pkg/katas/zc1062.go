// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:          "ZC1062",
		Title:       "Prefer `grep -E` over `egrep`",
		Description: "`egrep` is deprecated. Use `grep -E` instead.",
		Severity:    SeverityInfo,
		Check:       checkZC1062,
		Fix:         fixZC1062,
	})
}

// fixZC1062 rewrites `egrep` to `grep -E` at the command name
// position. Single replacement — arguments stay untouched.
func fixZC1062(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "egrep" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("egrep"),
		Replace: "grep -E",
	}}
}

func checkZC1062(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	if name, ok := cmd.Name.(*ast.Identifier); ok && name.Value == "egrep" {
		return []Violation{{
			KataID:  "ZC1062",
			Message: "`egrep` is deprecated. Use `grep -E` instead.",
			Line:    name.Token.Line,
			Column:  name.Token.Column,
			Level:   SeverityInfo,
		}}
	}

	return nil
}

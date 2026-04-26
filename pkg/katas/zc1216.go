// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1216",
		Title:    "Avoid `nslookup` — use `dig` or `host` for DNS queries",
		Severity: SeverityInfo,
		Description: "`nslookup` is deprecated in many distributions. " +
			"`dig` provides more detailed output and `host` is simpler for basic lookups.",
		Check: checkZC1216,
		Fix:   fixZC1216,
	})
}

// fixZC1216 rewrites `nslookup` to `host` at the command name position.
// `host <name>` matches the most common `nslookup <name>` invocation;
// arguments stay untouched and exotic nslookup-only flags need manual
// review.
func fixZC1216(node ast.Node, v Violation, _ []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "nslookup" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("nslookup"),
		Replace: "host",
	}}
}

func checkZC1216(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "nslookup" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1216",
		Message: "Avoid `nslookup` — it is deprecated on many systems. " +
			"Use `dig` for detailed DNS queries or `host` for simple lookups.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}

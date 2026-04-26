// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1191",
		Title:    "Avoid `clear` command — use ANSI escape sequences",
		Severity: SeverityStyle,
		Description: "`clear` spawns an external process for screen clearing. " +
			"Use `print -n '\\e[2J\\e[H'` for faster terminal clearing.",
		Check: checkZC1191,
		Fix:   fixZC1191,
	})
}

// fixZC1191 rewrites a bare `clear` identifier into the equivalent
// ANSI-escape `print` invocation, avoiding the external process. The
// `$'...'` quoting is required so the lexer interprets the escape
// codes; plain single quotes pass them through literally. The `-rn`
// flag-bundle matches the canonical `print -rn` form ZShellCheck
// recommends elsewhere (see ZC1017, ZC1118), so the rewrite is
// idempotent on re-run.
func fixZC1191(node ast.Node, v Violation, _ []byte) []FixEdit {
	ident, ok := node.(*ast.Identifier)
	if !ok || ident == nil || ident.Value != "clear" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("clear"),
		Replace: "print -rn $'\\e[2J\\e[H'",
	}}
}

func checkZC1191(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok || ident == nil || ident.Value != "clear" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1191",
		Message: "Use `print -n '\\e[2J\\e[H'` instead of `clear`. " +
			"ANSI escape sequences avoid spawning an external process.",
		Line:   ident.Token.Line,
		Column: ident.Token.Column,
		Level:  SeverityStyle,
	}}
}

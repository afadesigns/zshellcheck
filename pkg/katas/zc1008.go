// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.LetStatementNode, Kata{
		ID:    "ZC1008",
		Title: "Use `\\$(())` for arithmetic operations",
		Description: "The `let` command is a shell builtin, but the `\\$(())` syntax is more portable " +
			"and generally preferred for arithmetic operations in Zsh. It's also more powerful as it " +
			"can be used in more contexts.",
		Severity: SeverityStyle,
		Check:    checkZC1008,
		// Reuse ZC1013's `let NAME=EXPR` → `(( NAME = EXPR ))` rewrite.
		// All three of ZC1008, ZC1013, ZC1022 fire on the same `let`
		// shape and want the same arithmetic-command form; the
		// conflict resolver dedupes overlapping edits.
		Fix: fixZC1013,
	})
}

func checkZC1008(node ast.Node) []Violation {
	// Duplicate check for 'let' covered by ZC1013?
	// ZC1008 title says \$(()) which is expansion.
	// But check was for LetStatement.
	// Let's keep it as 'let' check for now to match original intent, maybe redundant.
	stmt, ok := node.(*ast.LetStatement)
	if !ok {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1008",
		Message: "Use `\\$(())` for arithmetic operations instead of `let`.",
		Line:    stmt.TokenLiteralNode().Line,
		Column:  stmt.TokenLiteralNode().Column,
		Level:   SeverityStyle,
	}}
}

// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.LetStatementNode, Kata{
		ID:    "ZC1013",
		Title: "Use `((...))` for arithmetic operations instead of `let`",
		Description: "The `let` command is a shell builtin, but the `((...))` syntax is more portable " +
			"and generally preferred for arithmetic operations in Zsh.",
		Severity: SeverityInfo,
		Check:    checkZC1013,
		Fix:      fixZC1013,
	})
}

// fixZC1013 rewrites `let NAME=EXPR` -> `(( NAME = EXPR ))`. The
// replacement spans from the `let` keyword to the end of the logical
// line (first `;`, `\n`, or EOF). The original expression text is
// preserved byte-identical so operator precedence and side effects
// match the source. Multi-assignment forms (`let a=1 b=2`) are not
// attempted — the Fix bails when the AST does not have a single
// Name/Value pair.
func fixZC1013(node ast.Node, v Violation, source []byte) []FixEdit {
	stmt, ok := node.(*ast.LetStatement)
	if !ok {
		return nil
	}
	if stmt.Name == nil || stmt.Value == nil {
		return nil
	}
	// Defer to ZC1032 when the value matches the C-style
	// increment/decrement shape so the rewrite produces the
	// idiomatic `(( i++ ))` / `(( i-- ))` form rather than the
	// generic `(( i = i+1 ))` form. Both fixes span the same source
	// range, so emitting both would lose ZC1032's narrower output to
	// the conflict resolver's deterministic tie-break.
	if _, increment := zc1032Op(stmt); increment {
		return nil
	}
	start := LineColToByteOffset(source, v.Line, v.Column)
	if start < 0 {
		return nil
	}
	// Scan for the end of the statement: first semicolon, newline,
	// or EOF.
	end := start
	for end < len(source) {
		b := source[end]
		if b == '\n' || b == ';' {
			break
		}
		end++
	}
	// `let NAME=EXPR` — after the keyword the source is NAME=EXPR.
	// Split on the first `=` to honour the original spelling (the
	// AST stores Name separately but Value.String() can lose
	// whitespace or parentheses that matter).
	prefix := "let "
	if start+len(prefix) > end {
		return nil
	}
	body := string(source[start+len(prefix) : end])
	eq := -1
	for i := 0; i < len(body); i++ {
		if body[i] == '=' {
			eq = i
			break
		}
	}
	if eq < 0 {
		return nil
	}
	name := body[:eq]
	rhs := body[eq+1:]
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  end - start,
		Replace: "(( " + name + " = " + rhs + " ))",
	}}
}

func checkZC1013(node ast.Node) []Violation {
	stmt, ok := node.(*ast.LetStatement)
	if !ok {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1013",
		Message: "Use `((...))` for arithmetic operations instead of `let`.",
		Line:    stmt.TokenLiteralNode().Line,
		Column:  stmt.TokenLiteralNode().Column,
		Level:   SeverityInfo,
	}}
}

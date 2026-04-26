// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.CommandSubstitutionNode, Kata{
		ID:    "ZC1002",
		Title: "Use $(...) instead of backticks",
		Description: "Backticks are the old-style command substitution. " +
			"$(...) is nesting-safe, easier to read, and generally preferred.",
		Severity: SeverityStyle,
		Check:    checkZC1002,
		Fix:      fixZC1002,
	})
}

// fixZC1002 rewrites “ `cmd` “ -> `$(cmd)`. The violation's Line and
// Column point at the opening backtick. Locate the matching closing
// backtick by scanning forward, respecting backslash escapes, and emit
// a single replacement edit spanning both delimiters. Unterminated
// backtick spans are skipped (the parser rejects them earlier; this is
// defensive).
func fixZC1002(node ast.Node, v Violation, source []byte) []FixEdit {
	cs, ok := node.(*ast.CommandSubstitution)
	if !ok {
		return nil
	}
	start := LineColToByteOffset(source, v.Line, v.Column)
	if start < 0 || start >= len(source) || source[start] != '`' {
		return nil
	}

	// Walk forward to find the matching closing backtick. Double-quoted
	// strings and escaped backticks inside the span are preserved.
	end := -1
	for i := start + 1; i < len(source); i++ {
		switch source[i] {
		case '\\':
			i++ // skip escaped char
		case '`':
			end = i
		}
		if end >= 0 {
			break
		}
	}
	if end < 0 {
		return nil
	}

	inner := string(source[start+1 : end])
	// Defensive: if the inner payload already contains `$(...)` shaped
	// parens we could still round-trip, but our parser produced
	// cs.Command so use its stringified form as a sanity cross-check.
	if cs.Command != nil && cs.Command.String() == "" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  end - start + 1,
		Replace: "$(" + inner + ")",
	}}
}

func checkZC1002(node ast.Node) []Violation {
	violations := []Violation{}

	if cs, ok := node.(*ast.CommandSubstitution); ok {
		violations = append(violations, Violation{
			KataID: "ZC1002",
			Message: "Use $(...) instead of backticks for command substitution. " +
				"The `$(...)` syntax is more readable and can be nested easily.",
			Line:   cs.Token.Line,
			Column: cs.Token.Column,
			Level:  SeverityStyle,
		})
	}

	return violations
}

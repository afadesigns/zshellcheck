// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.ForLoopStatementNode, Kata{
		ID:    "ZC1040",
		Title: "Use (N) nullglob qualifier for globs in loops",
		Description: "In Zsh, a glob that matches nothing (e.g., `*.txt`) will cause an error by default. " +
			"Use the `(N)` glob qualifier to make it null (empty) if no matches found, preventing the error.",
		Severity: SeverityStyle,
		Check:    checkZC1040,
		Fix:      fixZC1040,
	})
}

// fixZC1040 appends `(N)` after a glob pattern in a `for` loop item
// list, turning `for f in *.txt` into `for f in *.txt(N)` so an
// empty match produces an empty iterator instead of an error.
// Span scanning ends at the first unescaped whitespace / delimiter.
func fixZC1040(_ ast.Node, v Violation, source []byte) []FixEdit {
	start := LineColToByteOffset(source, v.Line, v.Column)
	if start < 0 || start >= len(source) {
		return nil
	}
	argLen := unquotedArgLen(source, start)
	if argLen == 0 {
		return nil
	}
	end := start + argLen
	endLine, endCol := offsetLineColZC1040(source, end)
	if endLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    endLine,
		Column:  endCol,
		Length:  0,
		Replace: "(N)",
	}}
}

func offsetLineColZC1040(source []byte, offset int) (int, int) {
	if offset < 0 || offset > len(source) {
		return -1, -1
	}
	line := 1
	col := 1
	for i := 0; i < offset; i++ {
		if source[i] == '\n' {
			line++
			col = 1
			continue
		}
		col++
	}
	return line, col
}

func checkZC1040(node ast.Node) []Violation {
	loop, ok := node.(*ast.ForLoopStatement)
	if !ok {
		return nil
	}

	// Only check "for i in items..." style loops, not arithmetic loops
	if loop.Items == nil {
		return nil
	}

	violations := []Violation{}

	for _, item := range loop.Items {
		// We are looking for string literals that look like globs (contain *, ?, etc)
		// but do NOT contain (N) or (N-...) qualifiers.

		val := getStringValue(item)

		// If it is quoted, it is NOT a glob expansion.
		if len(val) > 0 && (val[0] == '"' || val[0] == '\'') {
			continue
		}

		if isGlob(val) && !hasNullGlobQualifier(val) {
			violations = append(violations, Violation{
				KataID: "ZC1040",
				Message: "Glob pattern '" + val + "' may error if no files match. " +
					"Append '(N)' to enable nullglob behavior: '" + val + "(N)'",
				Line:   item.TokenLiteralNode().Line,
				Column: item.TokenLiteralNode().Column,
				Level:  SeverityStyle,
			})
		}
	}

	return violations
}

func getStringValue(node ast.Node) string {
	switch n := node.(type) {
	case *ast.StringLiteral:
		return n.Value
	case *ast.ConcatenatedExpression:
		var sb strings.Builder
		for _, p := range n.Parts {
			sb.WriteString(getStringValue(p))
		}
		return sb.String()
	case *ast.Identifier:
		return n.Value
	case *ast.GroupedExpression:
		return "(" + getStringValue(n.Expression) + ")"
	case *ast.ArrayLiteral:
		var sb strings.Builder
		sb.WriteString("(")
		for i, el := range n.Elements {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(getStringValue(el))
		}
		sb.WriteString(")")
		return sb.String()
	default:
		// Fallback for operators treated as literals (like *)
		return n.TokenLiteral()
	}
}

func isGlob(s string) bool {
	// Simple check for common glob characters
	return strings.ContainsAny(s, "*?[]")
}

func hasNullGlobQualifier(s string) bool {
	// Check for (N) at the end. Zsh qualifiers are at the end.
	// This is a naive check.
	return strings.Contains(s, "(N)") || strings.Contains(s, "(N") // (N) or (N...)
}

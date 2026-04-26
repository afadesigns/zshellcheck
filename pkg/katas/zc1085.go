// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.ForLoopStatementNode, Kata{
		ID:    "ZC1085",
		Title: "Quote variable expansions in `for` loops",
		Description: "Unquoted variable expansions in `for` loops are split by IFS (usually spaces). " +
			"This often leads to iterating over words instead of lines or array elements. Quote the expansion to preserve structure.",
		Severity: SeverityWarning,
		Check:    checkZC1085,
		Fix:      fixZC1085,
	})
}

// fixZC1085 wraps an unquoted expansion in a `for` loop item list
// with double-quotes. Two-edit insertion at the span start and end.
// Span uses the shared unquotedArgLen scanner so `${arr[@]}`,
// `$(cmd args)`, `${var:-default}` all stay whole.
func fixZC1085(_ ast.Node, v Violation, source []byte) []FixEdit {
	start := LineColToByteOffset(source, v.Line, v.Column)
	if start < 0 || start >= len(source) {
		return nil
	}
	argLen := unquotedArgLen(source, start)
	if argLen == 0 {
		return nil
	}
	endLine, endCol := offsetLineColZC1085(source, start+argLen)
	if endLine < 0 {
		return nil
	}
	return []FixEdit{
		{Line: v.Line, Column: v.Column, Length: 0, Replace: `"`},
		{Line: endLine, Column: endCol, Length: 0, Replace: `"`},
	}
}

func offsetLineColZC1085(source []byte, offset int) (int, int) {
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

func checkZC1085(node ast.Node) []Violation {
	loop, ok := node.(*ast.ForLoopStatement)
	if !ok {
		return nil
	}

	// If Items is nil or empty, it's either C-style or implicit `in "$@"`, ignore
	if len(loop.Items) == 0 {
		return nil
	}

	var violations []Violation

	for _, item := range loop.Items {
		if isUnquotedExpansion(item) {
			violations = append(violations, Violation{
				KataID:  "ZC1085",
				Message: "Unquoted variable expansion in for loop. This will split on IFS (usually space). Quote it to iterate over lines or array elements.",
				Line:    item.TokenLiteralNode().Line,
				Column:  item.TokenLiteralNode().Column,
				Level:   SeverityWarning,
			})
		}
	}

	return violations
}

func isUnquotedExpansion(expr ast.Expression) bool {
	// Check for Identifier (e.g. $var)
	if id, ok := expr.(*ast.Identifier); ok {
		return id.TokenLiteralNode().Type == "VARIABLE"
	}

	// Check for ArrayAccess (e.g. ${arr[@]})
	if _, ok := expr.(*ast.ArrayAccess); ok {
		return true
	}

	// Check for DollarParenExpression (e.g. $(cmd))
	if _, ok := expr.(*ast.DollarParenExpression); ok {
		return true
	}

	// Check for CommandSubstitution (e.g. `cmd`)
	if _, ok := expr.(*ast.CommandSubstitution); ok {
		return true
	}

	// Check for ConcatenatedExpression
	if concat, ok := expr.(*ast.ConcatenatedExpression); ok {
		inQuotes := false
		for _, part := range concat.Parts {
			if str, ok := part.(*ast.StringLiteral); ok {
				if str.Value == "\"" {
					inQuotes = !inQuotes
					continue
				}
				// Single quotes technically shouldn't appear here if parsed as StringLiteral?
			}

			if !inQuotes {
				if isUnquotedExpansion(part) {
					return true
				}
			}
		}
	}

	return false
}

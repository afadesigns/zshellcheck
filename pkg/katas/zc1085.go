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
		Check: checkZC1085,
	})
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
				if str.Value == "'" {
					// Single quotes technically shouldn't appear here if parsed as StringLiteral?
					// Parser consumes single quoted strings as single token usually.
					// But if we have mixed quoting...
					// For now, assume double quotes toggle.
				}
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

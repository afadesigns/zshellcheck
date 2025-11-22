package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.ForLoopStatementNode, Kata{
		ID:          "ZC1042",
		Title:       "Use \"$@\" to iterate over arguments",
		Description: "`$*` joins all arguments into a single string, which is rarely what you want in a loop. Use `\"$@\"` to iterate over each argument individually.",
		Check:       checkZC1042,
	})
}

func checkZC1042(node ast.Node) []Violation {
	loop, ok := node.(*ast.ForLoopStatement)
	if !ok {
		return nil
	}

	// Only check "for i in items..." style loops
	if loop.Items == nil {
		return nil
	}

	violations := []Violation{}

	for _, item := range loop.Items {
		// Check for "$*" (quoted) or $* (unquoted)
		
		// Helper to get raw value structure would be useful, but let's inspect manually.
		// 1. Unquoted $* -> Identifier with Value="$*"
		// 2. Quoted "$*" -> StringLiteral (if handled by lexer as one token) or ConcatenatedExpression?
		
		// In our parser/lexer, variables inside quotes often result in StringLiteral if simple, 
		// or if interpolated, we need to check the parts.
		// However, "$*" is special.
		
		found := false
		
		if ident, ok := item.(*ast.Identifier); ok {
			if ident.Value == "$*" {
				found = true
			}
		} else if str, ok := item.(*ast.StringLiteral); ok {
			// Check if it *contains* $* inside quotes.
			// Note: Our Lexer.readString now preserves quotes.
			// If Value is `"$"` that's bad.
			if strings.Contains(str.Value, "$*") {
				found = true
			}
		} else if concat, ok := item.(*ast.ConcatenatedExpression); ok {
			// Check parts for identifier $*
			for _, part := range concat.Parts {
				if ident, ok := part.(*ast.Identifier); ok && ident.Value == "$*" {
					found = true
					break
				}
				// Or string part containing it
				if str, ok := part.(*ast.StringLiteral); ok && strings.Contains(str.Value, "$*") {
					found = true
					break
				}
			}
		}

		if found {
			violations = append(violations, Violation{
				KataID:  "ZC1042",
				Message: "Use \"$@\" instead of \"$*\" (or $*) to iterate over arguments. \"$*\" merges arguments into a single string.",
				Line:    item.TokenLiteralNode().Line,
				Column:  item.TokenLiteralNode().Column,
			})
		}
	}

	return violations
}

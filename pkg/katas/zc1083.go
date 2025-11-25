package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.ConcatenatedExpressionNode, Kata{
		ID:    "ZC1083",
		Title: "Brace expansion limits cannot be variables",
		Description: "Brace expansion `{x..y}` happens before variable expansion. " +
			"`{1..$n}` will not work. Use `seq` or `for ((...))`.",
		Check: checkZC1083,
	})
	RegisterKata(ast.StringLiteralNode, Kata{
		ID:    "ZC1083",
		Title: "Brace expansion limits cannot be variables",
		Description: "Brace expansion `{x..y}` happens before variable expansion. " +
			"`{1..$n}` will not work. Use `seq` or `for ((...))`.",
		Check: checkZC1083,
	})
}

func checkZC1083(node ast.Node) []Violation {
	if strNode, ok := node.(*ast.StringLiteral); ok {
		val := strNode.Value
		if strings.Contains(val, "{") && strings.Contains(val, "..") && strings.Contains(val, "$") {
			return []Violation{{
				KataID:  "ZC1083",
				Message: "Brace expansion limits cannot be variables. `{...$var...}` is treated as a literal string. Use `seq` or `for ((...))`.",
				Line:    strNode.TokenLiteralNode().Line,
				Column:  strNode.TokenLiteralNode().Column,
			}}
		}
		return nil
	}

	concat, ok := node.(*ast.ConcatenatedExpression)
	if !ok {
		return nil
	}

	startIdx := -1
	var dotDotIndices []int
	var varIndices []int

	lastPartWasDot := false

	for i, part := range concat.Parts {
		if strNode, ok := part.(*ast.StringLiteral); ok {
			val := strNode.Value
			
			if strings.Contains(val, "{") {
				if startIdx == -1 {
					startIdx = i
				}
			}
			
			if strings.Contains(val, "..") {
				dotDotIndices = append(dotDotIndices, i)
				lastPartWasDot = false
			} else if val == "." {
				if lastPartWasDot {
					dotDotIndices = append(dotDotIndices, i-1) // Mark previous index as start of ..
					lastPartWasDot = false // Consumed
				} else {
					lastPartWasDot = true
				}
			} else {
				lastPartWasDot = false
			}
		} else {
			lastPartWasDot = false
			if _, ok := part.(*ast.IntegerLiteral); ok {
				continue
			}
			
			if idNode, ok := part.(*ast.Identifier); ok {
				if strings.Contains(idNode.Value, "..") {
					dotDotIndices = append(dotDotIndices, i)
				}
			}

			// Assume any other part is a variable or dynamic expansion
			varIndices = append(varIndices, i)
		}
	}

	if startIdx == -1 {
		return nil
	}

	// Check if .. is after startIdx
	hasDotDot := false
	for _, idx := range dotDotIndices {
		if idx > startIdx {
			hasDotDot = true
			break
		}
	}

	if !hasDotDot {
		return nil
	}

	// Check if variable is after startIdx
	hasVar := false
	for _, idx := range varIndices {
		if idx > startIdx {
			hasVar = true
			break
		}
	}

	if hasVar {
		return []Violation{{
			KataID:  "ZC1083",
			Message: "Brace expansion limits cannot be variables. `{...$var...}` is treated as a literal string. Use `seq` or `for ((...))`.",
			Line:    concat.TokenLiteralNode().Line,
			Column:  concat.TokenLiteralNode().Column,
		}}
	}

	return nil
}

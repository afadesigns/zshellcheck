package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.DoubleBracketExpressionNode, Kata{
		ID:    "ZC1079",
		Title: "Quote RHS of `==` in `[[ ... ]]` to prevent pattern matching",
		Description: "In `[[ ... ]]`, unquoted variable expansions on the right-hand side of `==` or `!=` " +
			"are treated as patterns (globbing). If you intend to compare strings literally, quote the variable.",
		Severity: SeverityWarning,
		Check:    checkZC1079,
		Fix:      fixZC1079,
	})
}

// fixZC1079 wraps an unquoted RHS variable reference inside `[[ … ]]`
// with double-quotes. Two edits: one `"` before the RHS token, one
// after. RHS span is measured from source so `${arr[$i]}` and
// `${var:-default}` stay whole. When the sibling LHS is an empty
// string literal, ZC1055's `-z` / `-n` rewrite takes priority and
// this fix no-ops to avoid overlapping edits.
func fixZC1079(node ast.Node, v Violation, source []byte) []FixEdit {
	if dbe, ok := node.(*ast.DoubleBracketExpression); ok {
		for _, el := range dbe.Elements {
			infix, ok := el.(*ast.InfixExpression)
			if !ok {
				continue
			}
			if infix.Operator != "==" && infix.Operator != "=" && infix.Operator != "!=" {
				continue
			}
			if isEmptyStringLiteral(infix.Left) || isEmptyStringLiteral(infix.Right) {
				// ZC1055 owns this rewrite; skip.
				return nil
			}
		}
	}
	start := LineColToByteOffset(source, v.Line, v.Column)
	if start < 0 || start >= len(source) {
		return nil
	}
	argLen := unquotedArgLen(source, start)
	if argLen == 0 {
		return nil
	}
	endOff := start + argLen
	endLine, endCol := offsetLineColZC1079(source, endOff)
	if endLine < 0 {
		return nil
	}
	return []FixEdit{
		{Line: v.Line, Column: v.Column, Length: 0, Replace: `"`},
		{Line: endLine, Column: endCol, Length: 0, Replace: `"`},
	}
}

func isEmptyStringLiteral(n ast.Node) bool {
	str, ok := n.(*ast.StringLiteral)
	if !ok {
		return false
	}
	return str.Value == `""` || str.Value == `''`
}

func offsetLineColZC1079(source []byte, offset int) (int, int) {
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

func checkZC1079(node ast.Node) []Violation {
	dbe, ok := node.(*ast.DoubleBracketExpression)
	if !ok {
		return nil
	}

	violations := []Violation{}

	for _, expr := range dbe.Elements {
		infix, ok := expr.(*ast.InfixExpression)
		if !ok {
			continue
		}

		// Check for equality/inequality operators
		if infix.Operator != "==" && infix.Operator != "=" && infix.Operator != "!=" {
			continue
		}

		// Check Right side
		// If it is an Identifier (variable), ArrayAccess, or Concatenated containing variable,
		// AND it is NOT quoted (not StringLiteral).

		// Note: Parser handles quoted strings as StringLiteral.
		// Unquoted $var is Identifier.

		isSuspicious := false
		var tokenNode ast.Node

		switch r := infix.Right.(type) {
		case *ast.Identifier:
			if len(r.Value) > 0 && r.Value[0] == '$' {
				isSuspicious = true
				tokenNode = r
			}
		case *ast.ArrayAccess:
			isSuspicious = true // ${arr[i]}
			tokenNode = r
		case *ast.InvalidArrayAccess:
			// ZC1001 covers syntax, but it's also unquoted.
			isSuspicious = true
			tokenNode = r
		case *ast.ConcatenatedExpression:
			// Check if any part is an unquoted variable
			for _, part := range r.Parts {
				if ident, ok := part.(*ast.Identifier); ok {
					if len(ident.Value) > 0 && ident.Value[0] == '$' {
						isSuspicious = true
						tokenNode = ident
						break
					}
				}
			}
		}

		if isSuspicious {
			violations = append(violations, Violation{
				KataID:  "ZC1079",
				Message: "Unquoted RHS matches as pattern. Quote to force string comparison: `\"$var\"`.",
				Line:    tokenNode.TokenLiteralNode().Line,
				Column:  tokenNode.TokenLiteralNode().Column,
				Level:   SeverityWarning,
			})
		}
	}

	return violations
}

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
		Check: checkZC1079,
	})
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
			})
		}
	}

	return violations
}

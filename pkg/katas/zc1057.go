package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1057",
		Title: "Avoid `ls` in assignments",
		Description: "Assigning the output of `ls` to a variable is fragile. " +
			"Use globs or arrays (e.g. `files=(*)`) to handle filenames correctly.",
		Check: checkZC1057,
	})
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:    "ZC1057",
		Title: "Avoid `ls` in assignments",
		Description: "Assigning the output of `ls` to a variable is fragile. " +
			"Use globs or arrays (e.g. `files=(*)`) to handle filenames correctly.",
		Check: checkZC1057,
	})
}

func checkZC1057(node ast.Node) []Violation {
	violations := []Violation{}

	checkAssignment := func(expr ast.Expression) {
		// Check if expr is an assignment to `ls` output
		// Usually ConcatenatedExpression: [Ident(var), String(=), DollarParen(ls)]
		if concat, ok := expr.(*ast.ConcatenatedExpression); ok {
			hasEquals := false
			for _, part := range concat.Parts {
				if str, ok := part.(*ast.StringLiteral); ok && str.Value == "=" {
					hasEquals = true
					continue
				}
				if hasEquals {
					// Check RHS for ls substitution
					if isLsSubstitution(part) {
						violations = append(violations, Violation{
							KataID:  "ZC1057",
							Message: "Avoid assigning `ls` output to a variable. Use globs (e.g. `files=(*)`) instead.",
							Line:    part.TokenLiteralNode().Line,
							Column:  part.TokenLiteralNode().Column,
						})
					}
				}
			}
		}
	}

	switch n := node.(type) {
	case *ast.SimpleCommand:
		checkAssignment(n.Name)
		for _, arg := range n.Arguments {
			checkAssignment(arg)
		}
	case *ast.InfixExpression:
		if n.Operator == "=" {
			if isLsSubstitution(n.Right) {
				violations = append(violations, Violation{
					KataID:  "ZC1057",
					Message: "Avoid assigning `ls` output to a variable. Use globs (e.g. `files=(*)`) instead.",
					Line:    n.TokenLiteralNode().Line,
					Column:  n.TokenLiteralNode().Column,
				})
			}
		}
	}

	return violations
}

func isLsSubstitution(node ast.Node) bool {
	// Reuse logic from ZC1050?
	// ZC1050 `getCommandFromSubstitution` returns the command.
	// We check if command is `ls`.

	var cmd ast.Expression

	switch n := node.(type) {
	case *ast.CommandSubstitution:
		cmd = n.Command
	case *ast.DollarParenExpression:
		cmd = n.Command
	default:
		return false
	}

	// Check if cmd is `ls`
	// cmd can be SimpleCommand or Infix (pipeline).
	// If simple command `ls ...`
	if simple, ok := cmd.(*ast.SimpleCommand); ok {
		if name, ok := simple.Name.(*ast.Identifier); ok && name.Value == "ls" {
			return true
		}
	}
	return false
}

package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1045",
		Title: "Declare and assign separately to avoid masking return values",
		Description: "Declaring a variable with `local var=$(cmd)` masks the return value of `cmd`. " +
			"The `local` command returns 0 (success) even if `cmd` fails. " +
			"Declare the variable first (`local var`), then assign it (`var=$(cmd)`).",
		Check: checkZC1045,
	})
}

func checkZC1045(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	name := cmd.Name.String()
	if name != "local" && name != "typeset" && name != "declare" && name != "readonly" {
		return nil
	}

	violations := []Violation{}

	for _, arg := range cmd.Arguments {
		// Check if arg is an assignment containing a command substitution
		if hasCommandSubstitutionAssignment(arg) {
			violations = append(violations, Violation{
				KataID: "ZC1045",
				Message: "Declare and assign separately to avoid masking return values. " +
					"`local var=$(cmd)` masks the exit code of `cmd`.",
				Line:   arg.TokenLiteralNode().Line,
				Column: arg.TokenLiteralNode().Column,
			})
		}
	}

	return violations
}

func hasCommandSubstitutionAssignment(arg ast.Expression) bool {
	// Argument structure depends on parsing.
	// Usually ConcatenatedExpression for `var=$(cmd)`: [Identifier(var), StringLiteral(=), DollarParenExpression]
	// Or `var=`cmd``: [Identifier(var), StringLiteral(=), CommandSubstitution]

	concat, ok := arg.(*ast.ConcatenatedExpression)
	if !ok {
		return false
	}

	hasEquals := false
	hasCmdSubst := false

	for _, part := range concat.Parts {
		if str, ok := part.(*ast.StringLiteral); ok && str.Value == "=" {
			hasEquals = true
			continue
		}

		if hasEquals {
			// Check if RHS has command substitution
			if isCommandSubstitution(part) {
				hasCmdSubst = true
			}
		}
	}

	return hasEquals && hasCmdSubst
}

func isCommandSubstitution(node ast.Node) bool {
	switch n := node.(type) {
	case *ast.CommandSubstitution:
		return true
	case *ast.DollarParenExpression:
		return true
	case *ast.ConcatenatedExpression:
		// Recursively check parts? e.g. `var="foo $(cmd)"`
		for _, p := range n.Parts {
			if isCommandSubstitution(p) {
				return true
			}
		}
	case *ast.StringLiteral:
		// Check for interpolation in double-quoted strings
		val := n.Value
		if len(val) >= 2 && val[0] == '"' && val[len(val)-1] == '"' {
			// Scan for $(...) or `...`
			// Simple heuristic: unescaped $ followed by ( or unescaped `
			for i := 0; i < len(val); i++ {
				if val[i] == '\\' {
					i++ // skip next
					continue
				}
				if val[i] == '`' {
					return true
				}
				if val[i] == '$' && i+1 < len(val) && val[i+1] == '(' {
					return true
				}
			}
		}
		return false
	}
	return false
}

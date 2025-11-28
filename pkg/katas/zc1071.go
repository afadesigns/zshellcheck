package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:          "ZC1071",
		Title:       "Use `+=` for appending to arrays",
		Description: "Appending to an array using `arr=($arr ...)` is verbose and slower. Use `arr+=(...)` instead.",
		Check:       checkZC1071,
	})
}

func checkZC1071(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	if len(cmd.Arguments) < 2 {
		return nil
	}

	// Check assignment operator
	assignOp := cmd.Arguments[0]
	if str, ok := assignOp.(*ast.StringLiteral); !ok || str.Value != "=" {
		return nil
	}

	// Check RHS
	rhs := cmd.Arguments[1]
	varName := cmd.Name.String()
	found := false

	// We only check if RHS is `GroupedExpression`.
	// If parser fails on `arr=($arr 4)`, we miss it.
	// But `arr=($arr)` works.
	// If parser supports `( ... )` as argument list in future, this will work.
	if grouped, ok := rhs.(*ast.GroupedExpression); ok {
		ast.Walk(grouped.Expression, func(n ast.Node) bool {
			// Check ArrayAccess (for ${var})
			if aa, ok := n.(*ast.ArrayAccess); ok {
				if ident, ok := aa.Left.(*ast.Identifier); ok && ident.Value == varName {
					found = true
					return false
				}
			}
			// Check Identifier with value "$var" or "${var}"
			if ident, ok := n.(*ast.Identifier); ok {
				if ident.Value == "$"+varName || ident.Value == "${"+varName+"}" {
					found = true
					return false
				}
			}
			// Check PrefixExpression like `$var`
			if prefix, ok := n.(*ast.PrefixExpression); ok && prefix.Operator == "$" {
				if ident, ok := prefix.Right.(*ast.Identifier); ok && ident.Value == varName {
					found = true
					return false
				}
			}
			return true
		})
	}

	if found {
		return []Violation{{
			KataID: "ZC1071",
			Message: "Appending to an array using `arr=($arr ...)` is verbose and slower. " +
				"Use `arr+=(...)` instead.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
		}}
	}

	return nil
}

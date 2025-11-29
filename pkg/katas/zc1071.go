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

	if len(cmd.Arguments) == 0 {
		return nil
	}

	varName := cmd.Name.String()
	var rhs ast.Expression

	arg0 := cmd.Arguments[0]

	if concat, ok := arg0.(*ast.ConcatenatedExpression); ok {
		if len(concat.Parts) >= 2 {
			if str, ok := concat.Parts[0].(*ast.StringLiteral); ok && str.Value == "=" {
				rhs = concat.Parts[1]
			}
		}
	} else if len(cmd.Arguments) >= 2 {
		if str, ok := arg0.(*ast.StringLiteral); ok && str.Value == "=" {
			rhs = cmd.Arguments[1]
		}
	}

	if rhs == nil {
		return nil
	}

	found := false

	checkNode := func(n ast.Node) bool {
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
		if prefix, ok := n.(*ast.PrefixExpression); ok {
			if prefix.Operator == "$" {
				if ident, ok := prefix.Right.(*ast.Identifier); ok && ident.Value == varName {
					found = true
					return false
				}
			}
		}
		return true
	}

	// Handle GroupedExpression (legacy/single element)
	if grouped, ok := rhs.(*ast.GroupedExpression); ok {
		ast.Walk(grouped.Expression, checkNode)
	}

	// Handle ArrayLiteral (multiple elements)
	if arrayLit, ok := rhs.(*ast.ArrayLiteral); ok {
		for _, elem := range arrayLit.Elements {
			if found {
				break
			}
			ast.Walk(elem, checkNode)
		}
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

package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.ForLoopStatementNode, Kata{
		ID:    "ZC1066",
		Title: "Avoid iterating over `cat` output",
		Description: "Iterating over `cat` output is fragile because lines can contain spaces. " +
			"Use `while IFS= read -r line; do ... done < file` or `($(<file))` array expansion.",
		Check: checkZC1066,
	})
}

func checkZC1066(node ast.Node) []Violation {
	loop, ok := node.(*ast.ForLoopStatement)
	if !ok {
		return nil
	}

	if loop.Items == nil {
		return nil
	}

	violations := []Violation{}

	for _, item := range loop.Items {
		// Check for $(cat ...) or `cat ...`
		// Reuse ZC1050 logic but for `cat`
		cmd := getCommandFromSubstitutionZC1066(item)
		if cmd != nil {
			if simpleCmd, ok := cmd.(*ast.SimpleCommand); ok {
				if name, ok := simpleCmd.Name.(*ast.Identifier); ok && name.Value == "cat" {
					violations = append(violations, Violation{
						KataID:  "ZC1066",
						Message: "Avoid iterating over `cat` output. Use `while read` loop or `($(<file))` for line-based iteration.",
						Line:    item.TokenLiteralNode().Line,
						Column:  item.TokenLiteralNode().Column,
					})
				}
			}
		}
	}

	return violations
}

func getCommandFromSubstitutionZC1066(node ast.Node) ast.Expression {
	switch n := node.(type) {
	case *ast.CommandSubstitution:
		return n.Command
	case *ast.DollarParenExpression:
		return n.Command
	case *ast.ConcatenatedExpression:
		// Check if any part is a substitution of cat
		for _, part := range n.Parts {
			if cmd := getCommandFromSubstitutionZC1066(part); cmd != nil {
				return cmd
			}
		}
	}
	return nil
}

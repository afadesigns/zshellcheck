package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.ForLoopStatementNode, Kata{
		ID:    "ZC1050",
		Title: "Avoid iterating over `ls` output",
		Description: "Iterating over `ls` output is fragile because filenames can contain spaces and newlines. " +
			"Use globs (e.g. `for f in *.txt`) instead.",
		Check: checkZC1050,
	})
}

func checkZC1050(node ast.Node) []Violation {
	loop, ok := node.(*ast.ForLoopStatement)
	if !ok {
		return nil
	}

	// Check loop items
	if loop.Items == nil {
		return nil
	}

	violations := []Violation{}

	for _, item := range loop.Items {
		// Check for $(ls ...) or `ls ...`
		cmd := getCommandFromSubstitution(item)
		if cmd != nil {
			if simpleCmd, ok := cmd.(*ast.SimpleCommand); ok {
				if name, ok := simpleCmd.Name.(*ast.Identifier); ok && name.Value == "ls" {
					violations = append(violations, Violation{
						KataID: "ZC1050",
						Message: "Avoid iterating over `ls` output. " +
							"Use globs (e.g. `*.txt`) to handle filenames with spaces correctly.",
						Line:   item.TokenLiteralNode().Line,
						Column: item.TokenLiteralNode().Column,
					})
				}
			}
		}
	}

	return violations
}

func getCommandFromSubstitution(node ast.Node) ast.Expression {
	switch n := node.(type) {
	case *ast.CommandSubstitution:
		return n.Command
	case *ast.DollarParenExpression:
		return n.Command
	case *ast.ConcatenatedExpression:
		// Check if any part is a substitution of ls
		for _, part := range n.Parts {
			if cmd := getCommandFromSubstitution(part); cmd != nil {
				return cmd
			}
		}
	}
	return nil
}

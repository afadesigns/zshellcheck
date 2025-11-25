package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.DoubleBracketExpressionNode, Kata{
		ID:    "ZC1091",
		Title: "Use `((...))` for arithmetic comparisons in `[[...]]`",
		Description: "The `[[ ... ]]` construct is primarily for string comparisons and file tests. " +
			"For arithmetic comparisons (`-eq`, `-lt`, etc.), use the dedicated arithmetic context `(( ... ))`. " +
			"It is cleaner and strictly numeric.",
		Check: checkZC1091,
	})
}

func checkZC1091(node ast.Node) []Violation {
	dbe, ok := node.(*ast.DoubleBracketExpression)
	if !ok {
		return nil
	}

	var violations []Violation

	visitor := func(n ast.Node) bool {
		if infix, ok := n.(*ast.InfixExpression); ok {
			switch infix.Operator {
			case "-eq", "-ne", "-lt", "-le", "-gt", "-ge":
				violations = append(violations, Violation{
					KataID:  "ZC1091",
					Message: "Use `(( ... ))` for arithmetic comparisons. e.g. `(( a < b ))` instead of `[[ a -lt b ]]`.",
					Line:    infix.TokenLiteralNode().Line,
					Column:  infix.TokenLiteralNode().Column,
				})
			}
		}
		return true
	}

	for _, expr := range dbe.Expressions {
		ast.Walk(expr, visitor)
	}

	return violations
}

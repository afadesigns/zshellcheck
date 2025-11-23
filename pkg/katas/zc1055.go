package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:    "ZC1055",
		Title: "Use `[[ -n/-z ]]` for empty string checks",
		Description: "Comparing with empty string is less idiomatic than using `[[ -z $var ]]` (is empty) " +
			"or `[[ -n $var ]]` (is not empty).",
		Check: checkZC1055,
	})
}

func checkZC1055(node ast.Node) []Violation {
	expr, ok := node.(*ast.InfixExpression)
	if !ok {
		return nil
	}

	// Check for == "" or != ""
	if expr.Operator != "==" && expr.Operator != "!=" {
		return nil
	}

	// Check if either side is empty string literal
	isEmptyString := func(n ast.Node) bool {
		if str, ok := n.(*ast.StringLiteral); ok {
			// Check for "" or ''
			val := str.Value
			return val == `""` || val == `''`
		}
		return false
	}

	if isEmptyString(expr.Left) || isEmptyString(expr.Right) {
		opSuggestion := "-z"
		if expr.Operator == "!=" {
			opSuggestion = "-n"
		}

		return []Violation{{
			KataID:  "ZC1055",
			Message: "Use `[[ " + opSuggestion + " ... ]]` instead of comparing with empty string.",
			Line:    expr.TokenLiteralNode().Line,
			Column:  expr.TokenLiteralNode().Column,
		}}
	}

	return nil
}

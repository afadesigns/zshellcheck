package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.ArithmeticCommandNode, Kata{
		ID:    "ZC1105",
		Title: "Avoid nested arithmetic expansions for clarity",
		Description: "While Zsh supports nested arithmetic expansions like `(( $((...)) ))`, " +
			"they can make code harder to read and reason about. Prefer flatter expressions " +
			"or temporary variables for intermediate results to improve clarity.",
		Check: checkZC1105,
	})
}

func checkZC1105(node ast.Node) []Violation {
	arithCmd, ok := node.(*ast.ArithmeticCommand)
	if !ok {
		return nil
	}

	// Check if the expression contains a nested arithmetic expansion
	// A simplified check: if the string representation contains another $(( or ((
	exprString := arithCmd.Expression.String()
	if strings.Contains(exprString, "$((") || strings.Contains(exprString, "((") {
		return []Violation{{
			KataID:  "ZC1105",
			Message: "Avoid nested arithmetic expansions. Use intermediate variables for clarity.",
			Line:    arithCmd.Token.Line,
			Column:  arithCmd.Token.Column,
		}}
	}

	return nil
}

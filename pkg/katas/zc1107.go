package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.DoubleBracketExpressionNode, Kata{
		ID:          "ZC1107",
		Title:       "Use (( ... )) for arithmetic conditions",
		Description: "Use `(( ... ))` for arithmetic comparisons instead of `[[ ... -gt ... ]]`. The double parenthesis syntax supports standard math operators (`>`, `<`, `==`, `!=`) and is optimized.",
		Severity:    Info,
		Check:       checkZC1107DoubleBracket,
	})

	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:          "ZC1107",
		Title:       "Use (( ... )) for arithmetic conditions",
		Description: "Use `(( ... ))` for arithmetic comparisons instead of `[ ... -eq ... ]`. The double parenthesis syntax supports standard math operators (`>`, `<`, `==`, `!=`) and is optimized.",
		Severity:    Info,
		Check:       checkZC1107SimpleCommand,
	})
}

func checkZC1107DoubleBracket(node ast.Node) []Violation {
	dbe := node.(*ast.DoubleBracketExpression)
	var violations []Violation

	// Helper to check infix expressions recursively
	check := func(n ast.Node) bool {
		if infix, ok := n.(*ast.InfixExpression); ok {
			switch infix.Operator {
			case "-eq", "-ne", "-lt", "-le", "-gt", "-ge":
				violations = append(violations, Violation{
					KataID:  "ZC1107",
					Message: "Prefer `(( ... ))` for arithmetic comparisons (e.g., `(( a > b ))`) over `[[ ... ]]` with flags like `" + infix.Operator + "`.",
					Line:    infix.TokenLiteralNode().Line,
					Column:  infix.TokenLiteralNode().Column,
				})
			}
		}
		return true
	}

	// Walk the elements of the double bracket expression
	for _, el := range dbe.Elements {
		ast.Walk(el, check)
	}

	return violations
}

func checkZC1107SimpleCommand(node ast.Node) []Violation {
	cmd := node.(*ast.SimpleCommand)

	// Check if command is "[" or "test"
	cmdName := cmd.Name.TokenLiteral()
	if cmdName != "[" && cmdName != "test" {
		return nil
	}

	var violations []Violation
	for _, arg := range cmd.Arguments {
		argText := arg.TokenLiteral()
		switch argText {
		case "-eq", "-ne", "-lt", "-le", "-gt", "-ge":
			violations = append(violations, Violation{
				KataID:  "ZC1107",
				Message: "Prefer `(( ... ))` for arithmetic comparisons (e.g., `(( a > b ))`) over `[ ... ]` with flags like `" + argText + "`.",
				Line:    arg.TokenLiteralNode().Line,
				Column:  arg.TokenLiteralNode().Column,
			})
		}
	}

	return violations
}

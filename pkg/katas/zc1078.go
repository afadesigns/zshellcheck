package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1078",
		Title: "Quote `$@` and `$*` when passing arguments",
		Description: "Using unquoted `$@` or `$*` splits arguments by IFS (usually space). " +
			"Use `\"$@\"` to preserve the original argument grouping, or `\"$*\"` to join them into a single string.",
		Check: checkZC1078,
	})
}

func checkZC1078(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	violations := []Violation{}

	for _, arg := range cmd.Arguments {
		// Check string representation to catch various parsed forms of $@ and $*
		// unquoted $@ might be parsed as Identifier "$@" -> String() == "$@"
		// unquoted $* might be parsed as GroupedExpression -> String() == "($*)"
		// or other variations depending on parser state (e.g. PrefixExpression)
		
		s := arg.String()
		
		// Removing parens from GroupedExpression string representation for checking
		// (Note: String() adds parens for GroupedExpression)
		if len(s) >= 2 && s[0] == '(' && s[len(s)-1] == ')' {
			s = s[1 : len(s)-1]
		}

		if s == "$@" || s == "$*" {
			violations = append(violations, Violation{
				KataID:  "ZC1078",
				Message: "Unquoted " + s + " splits arguments. Use \"" + s + "\" to preserve structure.",
				Line:    arg.TokenLiteralNode().Line,
				Column:  arg.TokenLiteralNode().Column,
			})
		}
	}

	return violations
}

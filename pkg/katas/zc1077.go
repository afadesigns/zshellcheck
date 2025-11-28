package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.CommandSubstitutionNode, Kata{
		ID:    "ZC1077",
		Title: "Prefer `${var:u/l}` over `tr` for case conversion",
		Description: "Using `tr` in a pipeline for simple case conversion is slower than using " +
			"Zsh's built-in parameter expansion flags `:u` (upper) and `:l` (lower).",
		Check: checkZC1077,
	})
	RegisterKata(ast.DollarParenExpressionNode, Kata{
		ID:    "ZC1077",
		Title: "Prefer `${var:u/l}` over `tr` for case conversion",
		Description: "Using `tr` in a pipeline for simple case conversion is slower than using " +
			"Zsh's built-in parameter expansion flags `:u` (upper) and `:l` (lower).",
		Check: checkZC1077,
	})
}

func checkZC1077(node ast.Node) []Violation {
	var command ast.Node

	switch n := node.(type) {
	case *ast.CommandSubstitution:
		command = n.Command
	case *ast.DollarParenExpression:
		command = n.Command
	default:
		return nil
	}

	// Check for pipeline: echo $var | tr ...
	infix, ok := command.(*ast.InfixExpression)
	if !ok || infix.Operator != "|" {
		return nil
	}

	// Right side must be `tr`
	rightCmd, ok := infix.Right.(*ast.SimpleCommand)
	if !ok || rightCmd.Name.String() != "tr" {
		return nil
	}

	// Check arguments of tr
	if len(rightCmd.Arguments) < 2 {
		return nil
	}

	arg1 := rightCmd.Arguments[0].String()
	arg2 := rightCmd.Arguments[1].String()

	// Check for 'a-z' 'A-Z' or similar patterns
	// We check both quoted and unquoted versions roughly
	isUpper := checkTrPattern(arg1, "a-z") && checkTrPattern(arg2, "A-Z")
	isLower := checkTrPattern(arg1, "A-Z") && checkTrPattern(arg2, "a-z")

	// Also check POSIX classes
	isUpperPosix := checkTrPattern(arg1, "[:lower:]") && checkTrPattern(arg2, "[:upper:]")
	isLowerPosix := checkTrPattern(arg1, "[:upper:]") && checkTrPattern(arg2, "[:lower:]")

	if isUpper || isUpperPosix {
		return []Violation{{
			KataID:  "ZC1077",
			Message: "Use `${var:u}` instead of `tr` for uppercase conversion. It is faster and built-in.",
			Line:    node.TokenLiteralNode().Line,
			Column:  node.TokenLiteralNode().Column,
		}}
	}

	if isLower || isLowerPosix {
		return []Violation{{
			KataID:  "ZC1077",
			Message: "Use `${var:l}` instead of `tr` for lowercase conversion. It is faster and built-in.",
			Line:    node.TokenLiteralNode().Line,
			Column:  node.TokenLiteralNode().Column,
		}}
	}

	return nil
}

func checkTrPattern(arg, pattern string) bool {
	// Remove quotes
	stripped := arg
	if len(arg) >= 2 && ((arg[0] == '"' && arg[len(arg)-1] == '"') || (arg[0] == '\'' && arg[len(arg)-1] == '\'')) {
		stripped = arg[1 : len(arg)-1]
	}

	// Simple containment check - robust enough for standard patterns
	// We check if the core pattern exists
	// e.g. 'a-z' matches "a-z", 'a-z', [a-z]

	// For strictness, let's just check substring
	return stripped == pattern || stripped == "["+pattern+"]"
}

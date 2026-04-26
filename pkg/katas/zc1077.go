// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
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
		Severity: SeverityStyle,
		Check:    checkZC1077,
	})
	RegisterKata(ast.DollarParenExpressionNode, Kata{
		ID:    "ZC1077",
		Title: "Prefer `${var:u/l}` over `tr` for case conversion",
		Description: "Using `tr` in a pipeline for simple case conversion is slower than using " +
			"Zsh's built-in parameter expansion flags `:u` (upper) and `:l` (lower).",
		Severity: SeverityStyle,
		Check:    checkZC1077,
	})
}

func checkZC1077(node ast.Node) []Violation {
	rightCmd, ok := zc1077TrPipeline(node)
	if !ok || len(rightCmd.Arguments) < 2 {
		return nil
	}
	arg1 := rightCmd.Arguments[0].String()
	arg2 := rightCmd.Arguments[1].String()

	if zc1077IsUpperPair(arg1, arg2) {
		return zc1077Hit(node, "u", "uppercase")
	}
	if zc1077IsLowerPair(arg1, arg2) {
		return zc1077Hit(node, "l", "lowercase")
	}
	return nil
}

// zc1077TrPipeline returns the right-hand `tr` command of a `cmd | tr`
// pipeline embedded in either a backtick or `$()` substitution.
func zc1077TrPipeline(node ast.Node) (*ast.SimpleCommand, bool) {
	var command ast.Node
	switch n := node.(type) {
	case *ast.CommandSubstitution:
		command = n.Command
	case *ast.DollarParenExpression:
		command = n.Command
	default:
		return nil, false
	}
	infix, ok := command.(*ast.InfixExpression)
	if !ok || infix.Operator != "|" {
		return nil, false
	}
	rightCmd, ok := infix.Right.(*ast.SimpleCommand)
	if !ok || rightCmd.Name.String() != "tr" {
		return nil, false
	}
	return rightCmd, true
}

func zc1077IsUpperPair(a, b string) bool {
	return (checkTrPattern(a, "a-z") && checkTrPattern(b, "A-Z")) ||
		(checkTrPattern(a, "[:lower:]") && checkTrPattern(b, "[:upper:]"))
}

func zc1077IsLowerPair(a, b string) bool {
	return (checkTrPattern(a, "A-Z") && checkTrPattern(b, "a-z")) ||
		(checkTrPattern(a, "[:upper:]") && checkTrPattern(b, "[:lower:]"))
}

func zc1077Hit(node ast.Node, flag, label string) []Violation {
	return []Violation{{
		KataID:  "ZC1077",
		Message: "Use `${var:" + flag + "}` instead of `tr` for " + label + " conversion. It is faster and built-in.",
		Line:    node.TokenLiteralNode().Line,
		Column:  node.TokenLiteralNode().Column,
		Level:   SeverityStyle,
	}}
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

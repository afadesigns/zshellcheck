// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:       "ZC1142",
		Title:    "Avoid chained `grep | grep` — combine patterns",
		Severity: SeverityStyle,
		Description: "Chaining `grep pattern1 | grep pattern2` spawns multiple processes. " +
			"Use `grep -E 'p1.*p2|p2.*p1'` or `awk` for multi-pattern matching.",
		Check: checkZC1142,
	})
}

func checkZC1142(node ast.Node) []Violation {
	pipe, ok := node.(*ast.InfixExpression)
	if !ok || pipe.Operator != "|" {
		return nil
	}

	leftCmd, ok := pipe.Left.(*ast.SimpleCommand)
	if !ok || !isCommandName(leftCmd, "grep") {
		return nil
	}

	rightCmd, ok := pipe.Right.(*ast.SimpleCommand)
	if !ok || !isCommandName(rightCmd, "grep") {
		return nil
	}

	return []Violation{{
		KataID: "ZC1142",
		Message: "Avoid chaining `grep | grep`. Combine into a single `grep -E` with alternation " +
			"or use `awk` for multi-pattern matching to reduce pipeline processes.",
		Line:   pipe.TokenLiteralNode().Line,
		Column: pipe.TokenLiteralNode().Column,
		Level:  SeverityStyle,
	}}
}

package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:       "ZC1190",
		Title:    "Combine chained `grep -v` into single invocation",
		Severity: SeverityStyle,
		Description: "`grep -v p1 | grep -v p2` spawns two processes. " +
			"Use `grep -v -e p1 -e p2` to combine exclusions in one invocation.",
		Check: checkZC1190,
	})
}

func checkZC1190(node ast.Node) []Violation {
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

	leftHasV := false
	rightHasV := false

	for _, arg := range leftCmd.Arguments {
		if arg.String() == "-v" {
			leftHasV = true
		}
	}
	for _, arg := range rightCmd.Arguments {
		if arg.String() == "-v" {
			rightHasV = true
		}
	}

	if leftHasV && rightHasV {
		return []Violation{{
			KataID: "ZC1190",
			Message: "Combine `grep -v p1 | grep -v p2` into `grep -v -e p1 -e p2`. " +
				"A single invocation avoids an unnecessary pipeline.",
			Line:   pipe.TokenLiteralNode().Line,
			Column: pipe.TokenLiteralNode().Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}

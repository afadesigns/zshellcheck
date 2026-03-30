package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:       "ZC1163",
		Title:    "Use `grep -m 1` instead of `grep | head -1`",
		Severity: SeverityStyle,
		Description: "`grep pattern | head -1` spawns two processes when `grep -m 1` does the same. " +
			"The `-m` flag stops after the first match, avoiding the pipeline.",
		Check: checkZC1163,
	})
}

func checkZC1163(node ast.Node) []Violation {
	pipe, ok := node.(*ast.InfixExpression)
	if !ok || pipe.Operator != "|" {
		return nil
	}

	grepCmd, ok := pipe.Left.(*ast.SimpleCommand)
	if !ok || !isCommandName(grepCmd, "grep") {
		return nil
	}

	headCmd, ok := pipe.Right.(*ast.SimpleCommand)
	if !ok || !isCommandName(headCmd, "head") {
		return nil
	}

	// Check head has -1 or -n 1
	for _, arg := range headCmd.Arguments {
		val := arg.String()
		if val == "-1" || val == "-n1" {
			return []Violation{{
				KataID: "ZC1163",
				Message: "Use `grep -m 1` instead of `grep | head -1`. " +
					"The `-m` flag stops after the first match without a pipeline.",
				Line:   pipe.TokenLiteralNode().Line,
				Column: pipe.TokenLiteralNode().Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

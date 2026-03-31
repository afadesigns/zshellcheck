package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:       "ZC1192",
		Title:    "Use `grep -c` instead of `grep | wc -l`",
		Severity: SeverityStyle,
		Description: "`grep pattern | wc -l` spawns two processes for counting matches. " +
			"Use `grep -c pattern` which counts internally without a pipeline.",
		Check: checkZC1192,
	})
}

func checkZC1192(node ast.Node) []Violation {
	pipe, ok := node.(*ast.InfixExpression)
	if !ok || pipe.Operator != "|" {
		return nil
	}

	grepCmd, ok := pipe.Left.(*ast.SimpleCommand)
	if !ok || !isCommandName(grepCmd, "grep") {
		return nil
	}

	wcCmd, ok := pipe.Right.(*ast.SimpleCommand)
	if !ok || !isCommandName(wcCmd, "wc") {
		return nil
	}

	for _, arg := range wcCmd.Arguments {
		if arg.String() == "-l" {
			return []Violation{{
				KataID: "ZC1192",
				Message: "Use `grep -c pattern` instead of `grep pattern | wc -l`. " +
					"The `-c` flag counts matches without a pipeline.",
				Line:   pipe.TokenLiteralNode().Line,
				Column: pipe.TokenLiteralNode().Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

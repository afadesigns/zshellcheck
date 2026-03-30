package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:       "ZC1146",
		Title:    "Avoid `cat file | awk` — pass file to awk directly",
		Severity: SeverityStyle,
		Description: "`cat file | awk` spawns an unnecessary cat process. " +
			"Pass the file directly as `awk '...' file`.",
		Check: checkZC1146,
	})
}

func checkZC1146(node ast.Node) []Violation {
	pipe, ok := node.(*ast.InfixExpression)
	if !ok || pipe.Operator != "|" {
		return nil
	}

	catCmd, ok := pipe.Left.(*ast.SimpleCommand)
	if !ok || !isCommandName(catCmd, "cat") {
		return nil
	}

	// cat must have exactly one file argument and no flags
	if len(catCmd.Arguments) != 1 {
		return nil
	}
	for _, arg := range catCmd.Arguments {
		if len(arg.String()) > 0 && arg.String()[0] == '-' {
			return nil
		}
	}

	// Right side must be awk/sed/sort or similar
	rightCmd, ok := pipe.Right.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	rightIdent, ok := rightCmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	name := rightIdent.Value
	if name == "awk" || name == "sed" || name == "sort" || name == "head" || name == "tail" {
		return []Violation{{
			KataID: "ZC1146",
			Message: "Pass the file directly to `" + name + "` instead of `cat file | " + name + "`. " +
				"Most text-processing tools accept file arguments.",
			Line:   pipe.TokenLiteralNode().Line,
			Column: pipe.TokenLiteralNode().Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}

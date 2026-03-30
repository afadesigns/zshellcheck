package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:    "ZC1093",
		Title: "Avoid useless `cat`",
		Description: "`cat file | command` spawns an unnecessary process. " +
			"Use `command < file` or pass the file as an argument directly.",
		Severity: SeverityStyle,
		Check:    checkZC1093,
	})
}

func checkZC1093(node ast.Node) []Violation {
	pipe, ok := node.(*ast.InfixExpression)
	if !ok || pipe.Operator != "|" {
		return nil
	}

	catCmd, ok := pipe.Left.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := catCmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cat" {
		return nil
	}

	// cat with flags like -n, -v, -e is not useless
	for _, arg := range catCmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] == '-' {
			return nil
		}
	}

	// Must have exactly one file argument
	if len(catCmd.Arguments) != 1 {
		return nil
	}

	return []Violation{{
		KataID: "ZC1093",
		Message: "`cat file | command` is inefficient. " +
			"Use `command < file` or pass the filename as an argument.",
		Line:   pipe.TokenLiteralNode().Line,
		Column: pipe.TokenLiteralNode().Column,
		Level:  SeverityStyle,
	}}
}

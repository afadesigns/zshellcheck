package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:    "ZC1126",
		Title: "Use `sort -u` instead of `sort | uniq`",
		Description: "`sort | uniq` spawns two processes when `sort -u` does the same in one. " +
			"Use `sort -u` to deduplicate sorted output efficiently.",
		Severity: SeverityStyle,
		Check:    checkZC1126,
	})
}

func checkZC1126(node ast.Node) []Violation {
	pipe, ok := node.(*ast.InfixExpression)
	if !ok || pipe.Operator != "|" {
		return nil
	}

	sortCmd, ok := pipe.Left.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if !isCommandName(sortCmd, "sort") {
		return nil
	}

	uniqCmd, ok := pipe.Right.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if !isCommandName(uniqCmd, "uniq") {
		return nil
	}

	// If uniq has flags like -c (count), -d (duplicates), skip
	for _, arg := range uniqCmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] == '-' {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1126",
		Message: "Use `sort -u` instead of `sort | uniq`. " +
			"Combining into one command avoids an unnecessary pipeline.",
		Line:   pipe.TokenLiteralNode().Line,
		Column: pipe.TokenLiteralNode().Column,
		Level:  SeverityStyle,
	}}
}

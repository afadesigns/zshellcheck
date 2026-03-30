package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:    "ZC1131",
		Title: "Avoid `cat file | while read` — use redirection",
		Description: "`cat file | while read line` spawns an unnecessary cat process " +
			"and runs the loop in a subshell. Use `while read line; do ...; done < file` instead.",
		Severity: SeverityStyle,
		Check:    checkZC1131,
	})
}

func checkZC1131(node ast.Node) []Violation {
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

	// cat must have exactly one file argument and no flags
	if len(catCmd.Arguments) != 1 {
		return nil
	}
	for _, arg := range catCmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] == '-' {
			return nil
		}
	}

	// Right side should involve 'while' or 'read'
	if right, ok := pipe.Right.(*ast.SimpleCommand); ok {
		rightIdent, ok := right.Name.(*ast.Identifier)
		if ok && rightIdent.Value == "read" {
			return []Violation{{
				KataID: "ZC1131",
				Message: "Use `while read line; do ...; done < file` instead of `cat file | while read line`. " +
					"Avoids unnecessary cat and subshell from the pipe.",
				Line:   pipe.TokenLiteralNode().Line,
				Column: pipe.TokenLiteralNode().Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

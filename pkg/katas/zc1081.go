package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:    "ZC1081",
		Title: "Use `${#var}` to get string length instead of `wc -c`",
		Description: "Using `echo $var | wc -c` involves a subshell and external command overhead. " +
			"Zsh has a built-in operator `${#var}` to get the length of a string instantly.",
		Check: checkZC1081,
	})
}

func checkZC1081(node ast.Node) []Violation {
	infix, ok := node.(*ast.InfixExpression)
	if !ok {
		return nil
	}

	if infix.Operator != "|" {
		return nil
	}

	// Check Right side: wc -c or wc -m
	rightCmd, ok := infix.Right.(*ast.SimpleCommand)
	if !ok || rightCmd.Name.String() != "wc" {
		return nil
	}

	isCharCount := false
	for _, arg := range rightCmd.Arguments {
		s := arg.String()
		if strings.Contains(s, "-c") || strings.Contains(s, "-m") {
			isCharCount = true
			break
		}
	}

	if !isCharCount {
		return nil
	}

	// Check Left side: echo ... or printf ...
	leftCmd, ok := infix.Left.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	cmdName := leftCmd.Name.String()
	if cmdName == "echo" || cmdName == "print" || cmdName == "printf" {
		return []Violation{{
			KataID:  "ZC1081",
			Message: "Use `${#var}` to get string length. Pipeline to `wc` is inefficient.",
			Line:    infix.TokenLiteralNode().Line,
			Column:  infix.TokenLiteralNode().Column,
		}}
	}

	return nil
}

package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:    "ZC1082",
		Title: "Prefer `${var//old/new}` over `sed` for simple replacements",
		Description: "Using `sed` for simple string replacement is slower than Zsh's built-in " +
			"parameter expansion. Use `${var/old/new}` (replace first) or `${var//old/new}` (replace all).",
		Check: checkZC1082,
	})
}

func checkZC1082(node ast.Node) []Violation {
	infix, ok := node.(*ast.InfixExpression)
	if !ok || infix.Operator != "|" {
		return nil
	}

	// Check Right side: sed
	rightCmd, ok := infix.Right.(*ast.SimpleCommand)
	if !ok || rightCmd.Name.String() != "sed" {
		return nil
	}

	// Check Left side: echo/printf/print
	leftCmd, ok := infix.Left.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	cmdName := leftCmd.Name.String()
	if cmdName != "echo" && cmdName != "print" && cmdName != "printf" {
		return nil
	}

	// Analyze sed arguments
	for _, arg := range rightCmd.Arguments {
		argStr := arg.String()
		// Remove quotes
		argStr = strings.Trim(argStr, "\"'")

		// Look for s/old/new/ or s/old/new/g
		// Basic check: starts with s/
		if strings.HasPrefix(argStr, "s/") || strings.HasPrefix(argStr, "s|") || strings.HasPrefix(argStr, "s@") {
			// It's a substitution
			return []Violation{{
				KataID:  "ZC1082",
				Message: "Use `${var//old/new}` for string replacement. Pipeline to `sed` is inefficient.",
				Line:    infix.TokenLiteralNode().Line,
				Column:  infix.TokenLiteralNode().Column,
			}}
		}
	}

	return nil
}

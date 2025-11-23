package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:          "ZC1041",
		Title:       "Do not use variables in printf format string",
		Description: "Using variables in `printf` format strings allows for format string attacks and unexpected behavior " +
			"if the variable contains `%`. Use `printf '%s' \"$var\"` instead.",
		Check:       checkZC1041,
	})
}

func checkZC1041(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	// Check if command is printf
	if cmdName, ok := cmd.Name.(*ast.Identifier); !ok || cmdName.Value != "printf" {
		return nil
	}

	if len(cmd.Arguments) == 0 {
		return nil
	}

	firstArg := cmd.Arguments[0]

	// The first argument should be a static StringLiteral.
	// If it is an Identifier ($var), ConcatenatedExpression ("$var"), or CommandSubstitution, warn.
	// Note: A StringLiteral might still contain interpolation if the lexer didn't split it,
	// but generally in this AST, StringLiteral is safe/static or single-quoted.
	// We warn if it's NOT a StringLiteral.

	_, isStringLiteral := firstArg.(*ast.StringLiteral)

	if !isStringLiteral {
		violations := []Violation{{
			KataID:  "ZC1041",
			Message: "Do not use variables in printf format string. Use `printf '..%s..' \"$var\"` instead.",
			Line:    firstArg.TokenLiteralNode().Line,
			Column:  firstArg.TokenLiteralNode().Column,
		}}
		return violations
	}

	// Even if it is a StringLiteral, it might be "$var" (interpolation).
	// We should inspect the value.
	if str, ok := firstArg.(*ast.StringLiteral); ok {
		val := str.Value
		// If it contains $ and is not single-quoted, it's likely a variable.
		// Heuristic: if it starts with " and contains $, it's risky.
		if len(val) > 0 && val[0] == '"' {
			// Check for $ not escaped? The lexer hands us the raw string usually.
			// Simple check: if it has unescaped $, flag it.
			// This is a basic heuristic.
			for i := 0; i < len(val); i++ {
				if val[i] == '$' && (i == 0 || val[i-1] != '\\') {
					violations := []Violation{{
						KataID:  "ZC1041",
						Message: "Do not use variables in printf format string. Use `printf '..%s..' \"$var\"` instead.",
						Line:    firstArg.TokenLiteralNode().Line,
						Column:  firstArg.TokenLiteralNode().Column,
					}}
					return violations
				}
			}
		}
	}

	return nil
}

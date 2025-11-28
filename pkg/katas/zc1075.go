package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1075",
		Title: "Quote variable expansions to prevent globbing",
		Description: "Unquoted variable expansions in Zsh are subject to globbing (filename generation). " +
			"If the variable contains characters like `*` or `?`, it might match files unexpectedly. " +
			"Use quotes `\"$var\"` to prevent this.",
		Check: checkZC1075,
	})
}

func checkZC1075(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	violations := []Violation{}

	for _, arg := range cmd.Arguments {
		// Check if argument is a simple identifier starting with $
		// or a braced expression ${...} that is NOT inside a string literal.
		// The parser might wrap these in different ways.

		// If it's a bare IdentifierNode (variable expansion), it's unquoted.
		if ident, ok := arg.(*ast.Identifier); ok {
			// Identifiers that start with $ are variable expansions
			if len(ident.Value) > 0 && ident.Value[0] == '$' {
				violations = append(violations, Violation{
					KataID:  "ZC1075",
					Message: "Unquoted variable expansion '" + ident.Value + "' is subject to globbing. Quote it: \"" + ident.Value + "\".",
					Line:    ident.Token.Line,
					Column:  ident.Token.Column,
				})
			}
		} else if _, ok := arg.(*ast.ArrayAccess); ok {
			// Array access ${arr[idx]} is also subject to globbing if unquoted
			violations = append(violations, Violation{
				KataID:  "ZC1075",
				Message: "Unquoted array access is subject to globbing. Quote it.",
				Line:    arg.TokenLiteralNode().Line,
				Column:  arg.TokenLiteralNode().Column,
			})
		} else if _, ok := arg.(*ast.InvalidArrayAccess); ok {
			// $arr[idx] - ZC1001 flags this, but it's also unquoted.
			// Let ZC1001 handle the syntax error, but ZC1075 could also flag globbing.
			// We'll skip to reduce noise.
		}

		// Note: StringLiteral arguments are quoted, so we don't check them.
		// But ConcatenatedExpression might contain unquoted parts.
		// e.g. $var/foo
		if concat, ok := arg.(*ast.ConcatenatedExpression); ok {
			for _, part := range concat.Parts {
				if ident, ok := part.(*ast.Identifier); ok {
					if len(ident.Value) > 0 && ident.Value[0] == '$' {
						violations = append(violations, Violation{
							KataID:  "ZC1075",
							Message: "Unquoted variable expansion '" + ident.Value + "' in concatenated string is subject to globbing.",
							Line:    ident.Token.Line,
							Column:  ident.Token.Column,
						})
					}
				}
			}
		}
	}

	return violations
}

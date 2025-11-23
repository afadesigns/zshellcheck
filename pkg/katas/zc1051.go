package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1051",
		Title: "Quote variables in `rm` to avoid globbing",
		Description: "`rm $VAR` is dangerous if `$VAR` contains spaces or glob characters. " +
			"Quote the variable (`rm \"$VAR\"`) to ensure safe deletion.",
		Check: checkZC1051,
	})
}

func checkZC1051(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	// Check if command is rm
	if name, ok := cmd.Name.(*ast.Identifier); !ok || name.Value != "rm" {
		return nil
	}

	violations := []Violation{}

	for _, arg := range cmd.Arguments {
		isUnquoted := false

		switch n := arg.(type) {
		case *ast.Identifier:
			// $VAR
			isUnquoted = true
		case *ast.PrefixExpression:
			// $var (if parsed as prefix)
			if n.Operator == "$" {
				isUnquoted = true
			}
		case *ast.ArrayAccess:
			// ${var[...]} unquoted
			// Zsh DOES NOT split unquoted variable expansions by default!
			// BUT it DOES glob them.
			// `rm $var`. If var="a b", it deletes "a b" (one file).
			// If var="*", it expands to all files.
			// So checking for globbing safety is key.
			// `rm \"$var\"` prevents globbing.
			isUnquoted = true
		case *ast.DollarParenExpression:
			// $(...)
			isUnquoted = true
		}

		if isUnquoted {
			violations = append(violations, Violation{
				KataID:  "ZC1051",
				Message: "Unquoted variable in `rm`. Quote it to prevent globbing (e.g. `rm \"$VAR\"`).",
				Line:    arg.TokenLiteralNode().Line,
				Column:  arg.TokenLiteralNode().Column,
			})
		}
	}

	return violations
}

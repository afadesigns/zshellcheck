package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1017",
		Title: "Use `print -r` to print strings literally",
		Description: "The `print` command interprets backslash escape sequences by default. " +
			"To print a string literally, use the `-r` option.",
		Check: checkZC1017,
	})
}

func checkZC1017(node ast.Node) []Violation {
	violations := []Violation{}

	if cmd, ok := node.(*ast.SimpleCommand); ok {
		if name, ok := cmd.Name.(*ast.Identifier); ok && name.Value == "print" {
			hasRFlag := false
			for _, arg := range cmd.Arguments {
				if prefix, ok := arg.(*ast.PrefixExpression); ok && prefix.Operator == "-" {
					if ident, ok := prefix.Right.(*ast.Identifier); ok && ident.Value == "r" {
						hasRFlag = true
						break
					}
				}
			}
			if !hasRFlag {
				violations = append(violations, Violation{
					KataID:  "ZC1017",
					Message: "Use `print -r` to print strings literally.",
					Line:    name.Token.Line,
					Column:  name.Token.Column,
				})
			}
		}
	}

	return violations
}

package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:          "ZC1003",
		Title:       "Use `((...))` for arithmetic comparisons instead of `[` or `test`",
		Description: "Bash/Zsh have a dedicated arithmetic context `((...))` " +
			"which is cleaner and faster than `[` or `test` for numeric comparisons.",
		Check: checkZC1003,
	})
}

func checkZC1003(node ast.Node) []Violation {
	violations := []Violation{}

	if cmd, ok := node.(*ast.SimpleCommand); ok {
		name := cmd.Name.String()
		if name == "[" || name == "test" {
			for _, arg := range cmd.Arguments {
				val := arg.String()
				// Trim parens added by AST String() method for expressions
				val = strings.Trim(val, "()")
				
				if val == "-eq" || val == "-ne" || val == "-lt" || val == "-le" || val == "-gt" || val == "-ge" {
					violations = append(violations, Violation{
						KataID:  "ZC1003",
						Message: "Use `((...))` for arithmetic comparisons instead of `[` or `test`.",
						Line:    cmd.Token.Line,
						Column:  cmd.Token.Column,
					})
					return violations
				}
			}
		}
	}

	return violations
}

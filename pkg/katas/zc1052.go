package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1052",
		Title: "Avoid `sed -i` for portability",
		Description: "`sed -i` usage varies between GNU/Linux and macOS/BSD. " +
			"macOS requires an extension argument (e.g. `sed -i ''`), while GNU does not. " +
			"Use a temporary file and `mv`, or `perl -i`, for portability.",
		Check: checkZC1052,
	})
}

func checkZC1052(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	if name, ok := cmd.Name.(*ast.Identifier); !ok || name.Value != "sed" {
		return nil
	}

	violations := []Violation{}

	for _, arg := range cmd.Arguments {
		argStr := arg.String()
		argStr = strings.Trim(argStr, "\"'")
		if strings.HasPrefix(argStr, "-") {
			if argStr == "-i" {
				violations = append(violations, Violation{
					KataID:  "ZC1052",
					Message: "`sed -i` is non-portable (GNU vs BSD differences). Use `perl -i` or a temporary file.",
					Line:    arg.TokenLiteralNode().Line,
					Column:  arg.TokenLiteralNode().Column,
				})
			}
		}
	}

	return violations
}

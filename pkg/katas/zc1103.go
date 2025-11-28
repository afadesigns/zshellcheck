package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:    "ZC1103",
		Title: "Suggest `path` array instead of `$PATH` string manipulation (direct assignment)",
		Description: "Zsh automatically maps the `$PATH` environment variable to the `$path` array. " +
			"Modifying `$path` is cleaner and less error-prone than manipulating the colon-separated `$PATH` string.",
		Check: checkZC1103,
	})
}

func checkZC1103(node ast.Node) []Violation {
	infixExp, ok := node.(*ast.InfixExpression)
	if !ok {
		return nil
	}

	if infixExp.Operator == "=" {
		if ident, ok := infixExp.Left.(*ast.Identifier); ok && ident.Value == "PATH" {
			// Check if the right-hand side is an old-style PATH manipulation
			if strings.Contains(infixExp.Right.String(), "$PATH") {
				return []Violation{{
					KataID:  "ZC1103",
					Message: "Use the `path` array instead of manually manipulating the `$PATH` string.",
					Line:    infixExp.Token.Line,
					Column:  infixExp.Token.Column,
				}}
			}
		}
	}

	return nil
}

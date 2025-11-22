package katas

import (

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IndexExpressionNode, Kata{
		ID:    "ZC1001",
		Title: "Use ${} for array element access",
		Description: "In Zsh, accessing array elements with `$my_array[1]` doesn't work as expected. " +
			"It tries to access an element from an array named `my_array[1]`. " +
			"The correct way to access an array element is to use `${my_array[1]}`.",
		Check: checkZC1001,
	})
	RegisterKata(ast.InvalidArrayAccessNode, Kata{
		ID:    "ZC1001",
		Title: "Use ${} for array element access",
		Description: "In Zsh, accessing array elements with `$my_array[1]` doesn't work as expected. " +
			"It tries to access an element from an array named `my_array[1]`. " +
			"The correct way to access an array element is to use `${my_array[1]}`.",
		Check: checkZC1001,
	})
}

func checkZC1001(node ast.Node) []Violation {
	violations := []Violation{}

	if indexExp, ok := node.(*ast.IndexExpression); ok {
		if ident, ok := indexExp.Left.(*ast.Identifier); ok {
			if len(ident.Value) > 0 && ident.Value[0] == '$' {
				violations = append(violations, Violation{
					KataID: "ZC1001",
					Message: "Use ${} for array element access. " +
						"Accessing array elements with `" + ident.Value + "[...]` is not the correct syntax in Zsh.",
					Line:   ident.Token.Line,
					Column: ident.Token.Column,
				})
			}
		}
	} else if arrayAccess, ok := node.(*ast.InvalidArrayAccess); ok {
		violations = append(violations, Violation{
			KataID: "ZC1001",
			Message: "Use ${} for array element access. " +
				"Accessing array elements with `$my_array[1]` is not the correct syntax in Zsh.",
			Line:   arrayAccess.Token.Line,
			Column: arrayAccess.Token.Column,
		})
	}

	return violations
}

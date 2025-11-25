package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.FunctionDefinitionNode, Kata{
		ID:    "ZC1086",
		Title: "Prefer `func() { ... }` over `function func { ... }`",
		Description: "The `function` keyword is optional in Zsh and non-standard in POSIX sh. " +
			"Using `func() { ... }` is more portable and consistent.",
		Check: checkZC1086,
	})
	RegisterKata(ast.FunctionLiteralNode, Kata{
		ID:    "ZC1086",
		Title: "Prefer `func() { ... }` over `function func { ... }`",
		Description: "The `function` keyword is optional in Zsh and non-standard in POSIX sh. " +
			"Using `func() { ... }` is more portable and consistent.",
		Check: checkZC1086,
	})
}

func checkZC1086(node ast.Node) []Violation {
	// Case 1: function my_func { ... } -> Parsed as FunctionLiteralNode
	if funcLit, ok := node.(*ast.FunctionLiteral); ok {
		if funcLit.TokenLiteral() == "function" {
			return []Violation{
				{
					KataID:  "ZC1086",
					Message: "Prefer `func() { ... }` over `function func { ... }` for portability and consistency.",
					Line:    funcLit.TokenLiteralNode().Line,
					Column:  funcLit.TokenLiteralNode().Column,
				},
			}
		}
	}

	// Case 2: my_func() { ... } -> Parsed as FunctionDefinitionNode
	if funcDef, ok := node.(*ast.FunctionDefinition); ok {
		if funcDef.TokenLiteral() == "function" {
			return []Violation{
				{
					KataID:  "ZC1086",
					Message: "Prefer `func() { ... }` over `function func { ... }` for portability and consistency.",
					Line:    funcDef.TokenLiteralNode().Line,
					Column:  funcDef.TokenLiteralNode().Column,
				},
			}
		}
	}

	return nil
}

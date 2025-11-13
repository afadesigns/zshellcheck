package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(Kata{
		ID:          "ZC1002",
		Title:       "Use $(...) instead of backticks for command substitution",
		Description: "The `$(...)` syntax is the modern, recommended way to perform command substitution. It is more readable and can be nested easily, unlike backticks.",
		Check:       checkZC1002,
	})
}

func checkZC1002(node ast.Node) []Violation {
	// This is a placeholder implementation.
	// The actual implementation will require modifications to the parser and AST.
	return nil
}

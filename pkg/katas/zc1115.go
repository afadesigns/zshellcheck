package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1115",
		Title: "Use Zsh string manipulation instead of `rev`",
		Description: "Zsh can reverse strings using parameter expansion. " +
			"Avoid spawning `rev` as an external process for simple string reversal.",
		Check: checkZC1115,
	})
}

func checkZC1115(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "rev" {
		return nil
	}

	// Only flag rev without file arguments (pipeline use)
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] != '-' {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1115",
		Message: "Use Zsh string manipulation instead of `rev`. " +
			"Parameter expansion can reverse strings without spawning an external process.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
	}}
}

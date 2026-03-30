package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1101",
		Title: "Use `$(( ))` instead of `bc` for simple arithmetic",
		Description: "Zsh supports arithmetic expansion with `$(( ))` and floating point via `zmodload zsh/mathfunc`. " +
			"Avoid piping to `bc` for simple calculations.",
		Severity: SeverityStyle,
		Check:    checkZC1101,
	})
}

func checkZC1101(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "bc" {
		return nil
	}

	// bc with file arguments is a valid external use
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] != '-' {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1101",
		Message: "Use `$(( ))` for arithmetic instead of `bc`. " +
			"Zsh arithmetic expansion avoids spawning an external process.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

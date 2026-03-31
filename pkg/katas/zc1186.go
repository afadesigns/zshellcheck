package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1186",
		Title:    "Use `unset -v` or `unset -f` for explicit unsetting",
		Severity: SeverityInfo,
		Description: "Bare `unset name` is ambiguous — it unsets variables first, then functions. " +
			"Use `unset -v` for variables or `unset -f` for functions to be explicit.",
		Check: checkZC1186,
	})
}

func checkZC1186(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "unset" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-v" || val == "-f" {
			return nil
		}
	}

	if len(cmd.Arguments) == 0 {
		return nil
	}

	return []Violation{{
		KataID: "ZC1186",
		Message: "Use `unset -v name` for variables or `unset -f name` for functions. " +
			"Bare `unset` is ambiguous about what is being removed.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}

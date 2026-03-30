package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1119",
		Title: "Use `$EPOCHSECONDS` instead of `date +%s`",
		Description: "Zsh provides `$EPOCHSECONDS` and `$EPOCHREALTIME` via `zsh/datetime` module. " +
			"Avoid spawning `date` for simple Unix timestamp retrieval.",
		Severity: SeverityStyle,
		Check:    checkZC1119,
	})
}

func checkZC1119(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "date" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := strings.Trim(arg.String(), "'\"")
		if val == "+%s" || val == "+%s%N" {
			return []Violation{{
				KataID: "ZC1119",
				Message: "Use `$EPOCHSECONDS` or `$EPOCHREALTIME` (via `zmodload zsh/datetime`) " +
					"instead of `date +%s`. Avoids spawning an external process.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

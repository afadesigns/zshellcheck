package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1399",
		Title:    "Use Zsh `$signals` array instead of `kill -l` for signal enumeration",
		Severity: SeverityStyle,
		Description: "Zsh exposes the `$signals` array (from `zsh/parameter`) holding all signal " +
			"names indexed from 0. `print -l $signals` produces the same list as `kill -l` " +
			"without spawning an external process.",
		Check: checkZC1399,
	})
}

func checkZC1399(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "kill" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-l" {
			return []Violation{{
				KataID: "ZC1399",
				Message: "Use Zsh `print -l $signals` (after `zmodload zsh/parameter`) instead " +
					"of `kill -l` for listing signal names.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

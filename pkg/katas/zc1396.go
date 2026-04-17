package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1396",
		Title:    "Avoid `unset -n` — Bash nameref semantics not in Zsh",
		Severity: SeverityError,
		Description: "Bash's `unset -n NAME` unsets the nameref itself rather than the target " +
			"variable it points to. Zsh does not implement namerefs; `unset -n` flags as an " +
			"error or unsets something unintended. Use `unset -v` for variable unset and " +
			"`unset -f` for function unset explicitly.",
		Check: checkZC1396,
	})
}

func checkZC1396(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "unset" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-n" {
			return []Violation{{
				KataID: "ZC1396",
				Message: "`unset -n` is a Bash nameref operation. Zsh does not honor it; use " +
					"`unset -v NAME` (variable) or `unset -f NAME` (function) explicitly.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}

	return nil
}

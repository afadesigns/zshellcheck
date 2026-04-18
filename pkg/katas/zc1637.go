package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1637",
		Title:    "Style: prefer Zsh `typeset -r NAME=value` over POSIX `readonly NAME=value`",
		Severity: SeverityStyle,
		Description: "Both `readonly NAME` and `typeset -r NAME` create a read-only parameter. " +
			"In Zsh the idiomatic form is `typeset -r` — it composes with other typeset flags " +
			"(`-ir` for readonly integer, `-xr` for readonly export, `-gr` to pin a readonly " +
			"global from inside a function). `readonly` works but reads as a Bash / POSIX-ism " +
			"in a Zsh codebase.",
		Check: checkZC1637,
	})
}

func checkZC1637(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "readonly" {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}

	return []Violation{{
		KataID: "ZC1637",
		Message: "`readonly` works but `typeset -r NAME=value` is the Zsh-native form and " +
			"composes with other typeset flags.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

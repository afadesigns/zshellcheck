package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1405",
		Title:    "Avoid `env -u VAR cmd` — use Zsh `(unset VAR; cmd)` subshell",
		Severity: SeverityStyle,
		Description: "`env -u VAR cmd` unsets a variable for a single command. In Zsh the " +
			"idiomatic form is a subshell: `(unset VAR; cmd)` — no external `env` spawn, and " +
			"the unset is naturally scoped to the subshell.",
		Check: checkZC1405,
	})
}

func checkZC1405(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "env" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-u" {
			return []Violation{{
				KataID: "ZC1405",
				Message: "Use `(unset VAR; cmd)` subshell instead of `env -u VAR cmd`. " +
					"In-shell scoping, no external process.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

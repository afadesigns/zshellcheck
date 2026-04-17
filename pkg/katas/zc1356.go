package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1356",
		Title:    "Use `read -A` instead of `read -a` for array read in Zsh",
		Severity: SeverityError,
		Description: "Zsh's `read` uses `-A` (uppercase A) to read into an array. Bash uses `-a` " +
			"(lowercase) for the same thing. In Zsh, `read -a` assigns a flag to a scalar " +
			"variable — not what Bash users expect. Use `-A` for portable-Zsh behavior.",
		Check: checkZC1356,
	})
}

func checkZC1356(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "read" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-a" {
			return []Violation{{
				KataID: "ZC1356",
				Message: "Use `read -A` (uppercase) in Zsh to read into an array. " +
					"`read -a` has different semantics in Zsh than in Bash.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}

	return nil
}

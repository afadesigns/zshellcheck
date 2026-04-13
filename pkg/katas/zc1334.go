package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1334",
		Title:    "Avoid `type -p` — use `whence -p` in Zsh",
		Severity: SeverityWarning,
		Description: "`type -p` is a Bash flag that prints the path of a command. " +
			"Zsh `type` does not support `-p`. Use `whence -p` to get " +
			"the path of an external command in Zsh.",
		Check: checkZC1334,
	})
}

func checkZC1334(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "type" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-p" || val == "-P" {
			return []Violation{{
				KataID:  "ZC1334",
				Message: "Avoid `type -p` in Zsh — use `whence -p` to get the command path instead.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
				Level:   SeverityWarning,
			}}
		}
	}

	return nil
}

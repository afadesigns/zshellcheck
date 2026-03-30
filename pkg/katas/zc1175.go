package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1175",
		Title:    "Avoid `tput` for simple ANSI colors — use Zsh `%F{color}`",
		Severity: SeverityStyle,
		Description: "Zsh prompt expansion `%F{red}` and `%f` handle colors natively. " +
			"Avoid spawning `tput` for simple color output in prompts and scripts.",
		Check: checkZC1175,
	})
}

func checkZC1175(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "tput" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "setaf" || val == "setab" || val == "sgr0" || val == "bold" {
			return []Violation{{
				KataID: "ZC1175",
				Message: "Use Zsh `%F{color}` / `%f` or `$fg[color]` / `$reset_color` instead of `tput`. " +
					"Zsh handles ANSI colors natively without spawning external processes.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

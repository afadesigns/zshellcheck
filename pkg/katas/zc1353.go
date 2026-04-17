package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1353",
		Title:    "Avoid `printf -v` — use `print -v` or command substitution in Zsh",
		Severity: SeverityStyle,
		Description: "`printf -v var fmt ...` is a Bash-ism. In Zsh use `print -v var -rf fmt ...` " +
			"or plain command substitution `var=$(printf fmt ...)`. `-v` is silently ignored by " +
			"POSIX printf, producing surprising bugs on portable scripts.",
		Check: checkZC1353,
	})
}

func checkZC1353(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "printf" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-v" {
			return []Violation{{
				KataID: "ZC1353",
				Message: "Avoid `printf -v` in Zsh — use `print -v var -rf fmt ...` or " +
					"`var=$(printf fmt ...)`. `-v` is Bash-specific and ignored elsewhere.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

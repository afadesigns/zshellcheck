package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1432",
		Title:    "Dangerous: `passwd -d user` — deletes the password, leaving the account passwordless",
		Severity: SeverityError,
		Description: "`passwd -d user` removes the password entirely, making the account usable " +
			"without any password (depending on PAM config). This is almost never what you want — " +
			"use `passwd -l user` to lock the account, or `usermod -L` + delete the ssh keys to " +
			"fully disable login.",
		Check: checkZC1432,
	})
}

func checkZC1432(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "passwd" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-d" {
			return []Violation{{
				KataID: "ZC1432",
				Message: "`passwd -d user` deletes the password — account becomes passwordless. " +
					"Use `passwd -l user` to lock, or `usermod -L` + delete SSH keys to disable login.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}

	return nil
}

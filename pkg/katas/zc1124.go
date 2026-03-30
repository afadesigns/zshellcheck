package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1124",
		Title: "Use `: > file` instead of `cat /dev/null > file` to truncate",
		Description: "Truncating a file with `cat /dev/null > file` spawns an unnecessary process. " +
			"Use `: > file` or simply `> file` in Zsh.",
		Check: checkZC1124,
	})
}

func checkZC1124(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cat" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "/dev/null" {
			return []Violation{{
				KataID: "ZC1124",
				Message: "Use `: > file` instead of `cat /dev/null > file` to truncate. " +
					"The `:` builtin avoids spawning cat.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
			}}
		}
	}

	return nil
}

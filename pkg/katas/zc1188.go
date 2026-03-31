package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1188",
		Title:    "Use Zsh `path+=()` instead of `export PATH=$PATH:dir`",
		Severity: SeverityStyle,
		Description: "Zsh ties the `path` array to `$PATH`. Use `path+=(dir)` to append " +
			"directories cleanly instead of string manipulation with `export PATH=`.",
		Check: checkZC1188,
	})
}

func checkZC1188(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "export" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 5 && val[:5] == "PATH=" {
			return []Violation{{
				KataID: "ZC1188",
				Message: "Use `path+=(dir)` instead of `export PATH=$PATH:dir`. " +
					"Zsh ties the `path` array to `$PATH` for cleaner manipulation.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1346",
		Title:    "Use Zsh `*(u:name:)` glob qualifier instead of `find -user`",
		Severity: SeverityStyle,
		Description: "Zsh's `*(u:name:)` and `*(u+uid)` glob qualifiers match files by owner " +
			"(name or numeric uid). The `*(U)` shorthand matches files owned by the current user. " +
			"Avoid `find -user` for the same selection.",
		Check: checkZC1346,
	})
}

func checkZC1346(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "find" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-user" || v == "-uid" || v == "-nouser" {
			return []Violation{{
				KataID: "ZC1346",
				Message: "Use Zsh `*(u:name:)` / `*(u+uid)` / `*(U)` glob qualifiers instead of " +
					"`find -user`/`-uid`/`-nouser`. Ownership predicates live entirely in the shell.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

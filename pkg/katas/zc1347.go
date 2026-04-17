package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1347",
		Title:    "Use Zsh `*(g:name:)` glob qualifier instead of `find -group`",
		Severity: SeverityStyle,
		Description: "Zsh's `*(g:name:)` and `*(g+gid)` glob qualifiers match files by group " +
			"(name or numeric gid). The `*(G)` shorthand matches files in the current user's group. " +
			"Avoid `find -group`/`-gid`/`-nogroup` for the same selection.",
		Check: checkZC1347,
	})
}

func checkZC1347(node ast.Node) []Violation {
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
		if v == "-group" || v == "-gid" || v == "-nogroup" {
			return []Violation{{
				KataID: "ZC1347",
				Message: "Use Zsh `*(g:name:)` / `*(g+gid)` / `*(G)` glob qualifiers instead of " +
					"`find -group`/`-gid`/`-nogroup`. Group predicates live entirely in the shell.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

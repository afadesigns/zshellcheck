package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1348",
		Title:    "Use Zsh glob type qualifiers instead of `find -type`",
		Severity: SeverityStyle,
		Description: "Zsh glob qualifiers select node type directly: `*(/)` directories, `*(.)` " +
			"regular files, `*(@)` symlinks, `*(=)` sockets, `*(p)` named pipes, `*(*)` " +
			"executable regular files, `*(%)` char/block devices. Avoid `find -type X` for " +
			"the same selection.",
		Check: checkZC1348,
	})
}

func checkZC1348(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "find" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-type" {
			return []Violation{{
				KataID: "ZC1348",
				Message: "Use Zsh glob type qualifiers (`*(/)`, `*(.)`, `*(@)`, `*(=)`, `*(p)`, `*(*)`, " +
					"`*(%)`) instead of `find -type`. No external process required.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

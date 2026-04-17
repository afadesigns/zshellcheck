package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1342",
		Title:    "Use Zsh `*(L0)` glob qualifier instead of `find -empty`",
		Severity: SeverityStyle,
		Description: "Zsh's `*(L0)` glob qualifier matches files with length 0. " +
			"Combine with `.` or `/` to restrict to regular files or directories. " +
			"Avoid shelling out to `find -empty` for the same result.",
		Check: checkZC1342,
	})
}

func checkZC1342(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "find" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-empty" {
			return []Violation{{
				KataID: "ZC1342",
				Message: "Use Zsh `*(L0)` glob qualifier instead of `find -empty`. " +
					"Add `.` for regular files only: `*(.L0)`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

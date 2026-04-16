package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1341",
		Title:    "Use Zsh `*(.x)` glob qualifier instead of `find -executable`",
		Severity: SeverityStyle,
		Description: "Zsh's `*(.x)` glob qualifier matches regular files that are executable. " +
			"Avoid shelling out to `find -executable` when the same selection is one glob away.",
		Check: checkZC1341,
	})
}

func checkZC1341(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "find" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-executable" {
			return []Violation{{
				KataID: "ZC1341",
				Message: "Use Zsh `*(.x)` glob qualifier instead of `find -executable`. " +
					"The `.` restricts to regular files and `x` to the executable bit.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

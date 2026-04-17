package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1344",
		Title:    "Use Zsh `*(L±Nk)` glob qualifier instead of `find -size`",
		Severity: SeverityStyle,
		Description: "Zsh's `*(LN)`, `*(L+N)`, `*(L-N)` match files by size in 512-byte blocks " +
			"(or bytes with a unit suffix: `k`, `m`, `p`). Same expressive power as " +
			"`find -size` without an external process.",
		Check: checkZC1344,
	})
}

func checkZC1344(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "find" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-size" {
			return []Violation{{
				KataID: "ZC1344",
				Message: "Use Zsh `*(L±N[kmp])` glob qualifier instead of `find -size`. " +
					"Unit suffixes: `k` kilobytes, `m` megabytes, `p` 512-byte blocks.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

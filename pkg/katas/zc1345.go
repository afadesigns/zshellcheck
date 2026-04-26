// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1345",
		Title:    "Use Zsh `*(f:mode:)` glob qualifier instead of `find -perm`",
		Severity: SeverityStyle,
		Description: "Zsh's `*(f:mode:)` glob qualifier matches files by permission mode. " +
			"Use octal (`*(f:0755:)`) or symbolic (`*(f:u+x:)`) inside the colon-delimited form. " +
			"Avoids spawning `find` for permission filters.",
		Check: checkZC1345,
	})
}

func checkZC1345(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "find" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-perm" {
			return []Violation{{
				KataID: "ZC1345",
				Message: "Use Zsh `*(f:mode:)` glob qualifier instead of `find -perm`. " +
					"Octal (`*(f:0755:)`) or symbolic (`*(f:u+x:)`) expressions are both supported.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1363",
		Title:    "Use Zsh `*(e:...:)` eval qualifier instead of `find -newer`/`-older`",
		Severity: SeverityStyle,
		Description: "Zsh's `*(e:expr:)` glob qualifier evaluates an arbitrary expression per match — " +
			"perfect for `-newer REF`-style predicates. Example: `*(e:'[[ $REPLY -nt reference ]]':)` " +
			"selects files newer than `reference`.",
		Check: checkZC1363,
	})
}

func checkZC1363(node ast.Node) []Violation {
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
		if v == "-newer" || v == "-anewer" || v == "-cnewer" ||
			v == "-neweraa" || v == "-newercm" || v == "-newermt" {
			return []Violation{{
				KataID: "ZC1363",
				Message: "Use Zsh `*(e:'[[ $REPLY -nt REF ]]':)` eval glob qualifier instead of " +
					"`find -newer`/`-anewer`/`-cnewer`/`-newerXY`. `$REPLY` holds the current match.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

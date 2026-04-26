// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1295",
		Title:    "Use `vared` instead of `read -e` for interactive editing in Zsh",
		Severity: SeverityStyle,
		Description: "Zsh provides `vared` for interactive editing of variables with full " +
			"ZLE support (tab completion, history, cursor movement). The `read -e` flag " +
			"is a Bash extension; Zsh `vared` is the native equivalent.",
		Check: checkZC1295,
	})
}

func checkZC1295(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "read" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-e" {
			return []Violation{{
				KataID:  "ZC1295",
				Message: "Use `vared` instead of `read -e` in Zsh. `vared` provides full ZLE editing support natively.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
				Level:   SeverityStyle,
			}}
		}
	}

	return nil
}

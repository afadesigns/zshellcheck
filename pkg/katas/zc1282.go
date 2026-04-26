// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1282",
		Title:    "Use Zsh `${var:r}` instead of `sed` to remove file extension",
		Severity: SeverityStyle,
		Description: "Zsh provides the `:r` modifier to remove a filename extension. " +
			"Using `sed` or `cut` to strip the extension is unnecessary when the built-in " +
			"parameter expansion handles it directly.",
		Check: checkZC1282,
	})
}

func checkZC1282(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "sed" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "'s/\\.[^.]*$//'" || val == "s/\\.[^.]*$//" ||
			val == "'s/\\.[^.]*$//g'" || val == "s/\\.[^.]*$//g" {
			return []Violation{{
				KataID:  "ZC1282",
				Message: "Use Zsh parameter expansion `${var:r}` to remove the file extension instead of `sed`.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
				Level:   SeverityStyle,
			}}
		}
	}

	return nil
}

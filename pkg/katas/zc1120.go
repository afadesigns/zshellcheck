// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1120",
		Title: "Use `$PWD` instead of `pwd`",
		Description: "Zsh maintains `$PWD` as a built-in variable tracking the current directory. " +
			"Avoid spawning `pwd` as an external process.",
		Severity: SeverityStyle,
		Check:    checkZC1120,
	})
}

func checkZC1120(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "pwd" {
		return nil
	}

	// pwd -P (physical) resolves symlinks — $PWD may not
	// Only flag bare pwd or pwd -L (logical, same as $PWD)
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-P" {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1120",
		Message: "Use `$PWD` instead of `pwd`. " +
			"Zsh maintains `$PWD` as a built-in variable, avoiding an external process.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

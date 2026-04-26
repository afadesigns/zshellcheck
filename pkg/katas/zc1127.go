// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1127",
		Title: "Avoid `ls` for counting files",
		Description: "Using `ls | wc -l` to count files spawns unnecessary processes. " +
			"Use Zsh glob qualifiers: `files=(*(N)); echo ${#files}` for file counting.",
		Severity: SeverityStyle,
		Check:    checkZC1127,
	})
}

func checkZC1127(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ls" {
		return nil
	}

	// Flag ls -1 (single column listing, typically used for counting)
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-1" {
			return []Violation{{
				KataID: "ZC1127",
				Message: "Use Zsh glob qualifiers `files=(*(N)); echo ${#files}` instead of `ls -1 | wc -l`. " +
					"Avoids spawning external processes for file counting.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

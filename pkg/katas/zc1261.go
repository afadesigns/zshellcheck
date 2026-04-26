// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1261",
		Title:    "Avoid piping `base64 -d` output to shell execution",
		Severity: SeverityError,
		Description: "Decoding base64 and piping to `sh`/`zsh`/`eval` is a code injection risk. " +
			"Always inspect decoded content before execution.",
		Check: checkZC1261,
	})
}

func checkZC1261(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "base64" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-d" || val == "-D" {
			return []Violation{{
				KataID: "ZC1261",
				Message: "Inspect `base64 -d` output before piping to execution. " +
					"Blindly executing decoded content is a code injection vector.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}

	return nil
}

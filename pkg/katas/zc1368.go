// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1368",
		Title:    "Avoid `sh -c` / `bash -c` inside a Zsh script — inline or use a function",
		Severity: SeverityStyle,
		Description: "Invoking `sh -c` or `bash -c` inside a Zsh script spawns a second shell, " +
			"loses access to the parent script's functions, arrays, and associative arrays, and " +
			"re-interprets POSIX-only syntax. Inline the code as a function or use `zsh -c` when " +
			"a subshell is truly required.",
		Check: checkZC1368,
	})
}

func checkZC1368(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "sh" && ident.Value != "bash" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-c" {
			return []Violation{{
				KataID: "ZC1368",
				Message: "Avoid `sh -c`/`bash -c` inside a Zsh script — inline the code as a " +
					"function to keep access to arrays, associative arrays, and Zsh features. " +
					"Use `zsh -c` only when a fresh shell is truly required.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1137",
		Title: "Avoid hardcoded `/tmp` paths",
		Description: "Hardcoded `/tmp` paths are predictable and may cause race conditions " +
			"or symlink attacks. Use `mktemp` or Zsh `=(...)` for safe temp files.",
		Severity: SeverityStyle,
		Check:    checkZC1137,
	})
}

func checkZC1137(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	// Skip mktemp itself (it creates temp files properly)
	cmdIdent, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if cmdIdent.Value == "mktemp" || cmdIdent.Value == "cd" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		// Match /tmp/something with a predictable name
		if len(val) > 5 && val[:5] == "/tmp/" && val != "/tmp" {
			// Skip if it uses a variable (dynamic path)
			hasVar := false
			for _, ch := range val {
				if ch == '$' {
					hasVar = true
					break
				}
			}
			if !hasVar {
				return []Violation{{
					KataID: "ZC1137",
					Message: "Avoid hardcoded `/tmp/` paths. Use `mktemp` or Zsh `=(cmd)` " +
						"for temp files to prevent race conditions and symlink attacks.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityStyle,
				}}
			}
		}
	}

	return nil
}

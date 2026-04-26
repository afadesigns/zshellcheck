// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1181",
		Title:    "Avoid `xdg-open`/`open` — use `$BROWSER` for portability",
		Severity: SeverityInfo,
		Description: "`xdg-open` is Linux-only, `open` is macOS-only. " +
			"Use `$BROWSER` or check `$OSTYPE` for cross-platform URL/file opening.",
		Check: checkZC1181,
	})
}

func checkZC1181(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "xdg-open" && ident.Value != "open" {
		return nil
	}

	if ident.Value == "open" {
		// open with flags like -a, -e is macOS-specific file opening
		for _, arg := range cmd.Arguments {
			val := arg.String()
			if len(val) > 0 && val[0] == '-' {
				return nil // Likely intentional macOS usage
			}
		}
	}

	return []Violation{{
		KataID: "ZC1181",
		Message: "Use `$BROWSER` or check `$OSTYPE` instead of `" + ident.Value + "` for portable " +
			"URL/file opening across Linux and macOS.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}

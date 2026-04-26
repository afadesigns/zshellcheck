// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1419",
		Title:    "Avoid `chmod 777` — grants world-writable access",
		Severity: SeverityWarning,
		Description: "Mode 777 (or 0777) grants read/write/execute to owner, group, and world. " +
			"Files become world-writable, which on a multi-user system or inside a container " +
			"with mapped UIDs is almost always wrong. Use 755 for executables, 644 for regular " +
			"files, 700 for private directories, or `umask`-aware helpers.",
		Check: checkZC1419,
	})
}

func checkZC1419(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "chmod" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "777" || v == "0777" || v == "a+rwx" || v == "ugo+rwx" {
			return []Violation{{
				KataID: "ZC1419",
				Message: "Avoid `chmod 777`/`a+rwx` — grants world-writable access. Prefer " +
					"restrictive modes (755, 644, 700, 600) matched to the actual file purpose.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}

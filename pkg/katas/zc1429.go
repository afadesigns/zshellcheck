// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1429",
		Title:    "Avoid `umount -f` / `-l` — force/lazy unmount masks real issues",
		Severity: SeverityWarning,
		Description: "`umount -f` forces the unmount even if the FS is busy; `-l` (lazy) " +
			"detaches immediately but keeps the FS in-use. Both can leave stale file handles " +
			"and data loss. Fix the underlying 'target busy' (use `lsof` / `fuser -m` to find " +
			"users) instead of forcing.",
		Check: checkZC1429,
	})
}

func checkZC1429(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "umount" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-f" || v == "-l" || v == "-fl" || v == "-lf" {
			return []Violation{{
				KataID: "ZC1429",
				Message: "`umount -f`/`-l` force/lazy unmount masks the underlying 'busy' error. " +
					"Find open files with `lsof` / `fuser -m` and close them properly.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}

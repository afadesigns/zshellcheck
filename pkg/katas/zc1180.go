// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1180",
		Title:    "Avoid `pgrep` for own background jobs — use Zsh job control",
		Severity: SeverityInfo,
		Description: "For managing your own background jobs, use Zsh job control (`jobs`, `kill %N`, " +
			"`fg`, `bg`) instead of `pgrep`/`pkill` which search system-wide.",
		Check: checkZC1180,
	})
}

func checkZC1180(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "pgrep" && ident.Value != "pkill" {
		return nil
	}

	// Only flag simple pgrep/pkill without complex flags
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 1 && val[0] == '-' && val != "-f" {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1180",
		Message: "For own background jobs, use Zsh job control (`jobs`, `kill %N`) instead of `" +
			ident.Value + "`. Job control is more precise for script-spawned processes.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}

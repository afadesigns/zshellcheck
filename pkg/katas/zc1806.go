// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1806",
		Title:    "Warn on `zmv 'PAT' 'REP'` without `-n` / `-i` — silent bulk rename",
		Severity: SeverityWarning,
		Description: "`zmv` (autoloaded from Zsh's functions) rewrites every filename that matches " +
			"the pattern in one shot. A small typo in the source pattern or replacement — " +
			"`*.jpg` vs `*.JPG`, a misplaced `(..)`, forgetting `**` recursion — can collide " +
			"names and silently overwrite files, since `zmv` aborts the batch only on its " +
			"own conflict check, not on semantic errors. Use `zmv -n 'PAT' 'REP'` first to " +
			"see the rename list, or `zmv -i` to prompt per file. Only drop the guard once " +
			"the preview matches what you expect.",
		Check: checkZC1806,
	})
}

func checkZC1806(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "zmv" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-n" || v == "--dry-run" || v == "-i" || v == "--interactive" {
			return nil
		}
		// combined short flags like `-Mn` or `-wn`.
		if len(v) > 1 && v[0] == '-' && v[1] != '-' {
			for _, c := range v[1:] {
				if c == 'n' || c == 'i' {
					return nil
				}
			}
		}
	}
	return []Violation{{
		KataID: "ZC1806",
		Message: "`zmv` without `-n` (dry-run) or `-i` (interactive) renames every " +
			"matched file in one shot — a pattern typo can collide names. Preview " +
			"with `zmv -n`, then re-run once the list looks right.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

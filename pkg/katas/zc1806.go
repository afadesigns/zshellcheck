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
	if !ok || CommandIdentifier(cmd) != "zmv" || len(cmd.Arguments) < 2 {
		return nil
	}
	for _, arg := range cmd.Arguments {
		if zc1806HasGuardFlag(arg.String()) {
			return nil
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

func zc1806HasGuardFlag(v string) bool {
	switch v {
	case "-n", "--dry-run", "-i", "--interactive":
		return true
	}
	if len(v) <= 1 || v[0] != '-' || v[1] == '-' {
		return false
	}
	for _, c := range v[1:] {
		if c == 'n' || c == 'i' {
			return true
		}
	}
	return false
}

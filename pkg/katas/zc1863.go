// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1863",
		Title:    "Warn on `unsetopt CASE_GLOB` — globs silently go case-insensitive across the script",
		Severity: SeverityWarning,
		Description: "`CASE_GLOB` on is the Zsh default: `*.log` matches `app.log` but not " +
			"`APP.LOG`, `[A-Z]*` is a real case-sensitive range, and `[[ $f == Foo* ]]` " +
			"keeps the distinction between `Foo1` and `foo1`. Turning it off (or " +
			"equivalently `setopt NO_CASE_GLOB`) silently re-evaluates every subsequent " +
			"pattern case-insensitively — `rm *.log` now sweeps `APP.LOG` up, pattern " +
			"dispatchers that used to distinguish `README` from `readme` stop doing so, " +
			"and hash maps keyed on glob-built labels start colliding. Keep the option " +
			"on at script level; request case-folding per-pattern with the Zsh qualifier " +
			"`(#i)*.log`.",
		Check: checkZC1863,
	})
}

func checkZC1863(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			if zc1863IsCaseGlob(arg.String()) {
				return zc1863Hit(cmd, "unsetopt "+arg.String())
			}
		}
	case "setopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOCASEGLOB" {
				return zc1863Hit(cmd, "setopt "+v)
			}
		}
	}
	return nil
}

func zc1863IsCaseGlob(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "CASEGLOB"
}

func zc1863Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1863",
		Message: "`" + where + "` flips every later glob to case-insensitive — " +
			"`rm *.log` sweeps `APP.LOG`, dispatchers keyed on case collisions. " +
			"Keep the option on; use `(#i)pattern` per-glob when you need " +
			"case-folding.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

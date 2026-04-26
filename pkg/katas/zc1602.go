// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1602",
		Title:    "Warn on `setopt KSH_ARRAYS` / `SH_WORD_SPLIT` — flips Zsh core semantics shell-wide",
		Severity: SeverityWarning,
		Description: "`KSH_ARRAYS` makes arrays 0-indexed (the Bash / ksh convention), breaking " +
			"every Zsh access that uses `[1]` for the first element. `SH_WORD_SPLIT` makes " +
			"unquoted `$var` word-split on `IFS`, breaking the core Zsh promise that `echo " +
			"$x` passes exactly one argument. Setting either globally is a bug-magnet — pre-" +
			"existing code silently misbehaves from that line on. If you need the semantics " +
			"only inside a function, scope it with `emulate -L ksh` or `emulate -L sh`.",
		Check: checkZC1602,
	})
}

func checkZC1602(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "setopt" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		raw := arg.String()
		norm := strings.ToLower(strings.ReplaceAll(raw, "_", ""))
		if norm == "ksharrays" || norm == "shwordsplit" {
			return []Violation{{
				KataID: "ZC1602",
				Message: "`setopt " + raw + "` flips Zsh core semantics for the whole shell " +
					"— pre-existing code silently misbehaves. Scope with `emulate -L ksh` / " +
					"`emulate -L sh` inside the function that needs the mode.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

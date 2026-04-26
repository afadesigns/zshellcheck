// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1895",
		Title:    "Warn on `setopt NUMERIC_GLOB_SORT` — glob output switches from lexicographic to numeric order",
		Severity: SeverityWarning,
		Description: "`NUMERIC_GLOB_SORT` is off by default: `ls *.log` returns filenames in the " +
			"collation order the filesystem-iteration/sort step produces (lexicographic " +
			"in the C locale, so `app-1.log`, `app-10.log`, `app-2.log`). Turning it on " +
			"makes every subsequent glob and array expansion sort numeric runs " +
			"numerically — the same glob now returns `app-1.log`, `app-2.log`, " +
			"`app-10.log`. Scripts that tail the \"latest\" file by taking the last array " +
			"element, pipelines that expect a specific stable order, and backup rotations " +
			"built on `*[0-9].tar` silently shuffle. Keep the option off script-wide; " +
			"request numeric sort per-glob with the `*(n)` qualifier when needed.",
		Check: checkZC1895,
	})
}

func checkZC1895(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "setopt":
		for _, arg := range cmd.Arguments {
			if zc1895IsNumericGlobSort(arg.String()) {
				return zc1895Hit(cmd, "setopt "+arg.String())
			}
		}
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NONUMERICGLOBSORT" {
				return zc1895Hit(cmd, "unsetopt "+v)
			}
		}
	}
	return nil
}

func zc1895IsNumericGlobSort(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "NUMERICGLOBSORT"
}

func zc1895Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1895",
		Message: "`" + where + "` switches every later glob to numeric sort — log " +
			"rotations sorted on numeric suffixes silently shuffle. Keep it off; " +
			"use the per-glob `*(n)` qualifier when needed.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

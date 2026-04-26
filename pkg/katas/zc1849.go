// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1849",
		Title:    "Warn on `setopt ALL_EXPORT` — every later `var=value` silently becomes `export var=value`",
		Severity: SeverityWarning,
		Description: "`ALL_EXPORT` (POSIX `set -a` equivalent, off by default) tells Zsh to mark " +
			"every parameter assignment for export as soon as it is created, so " +
			"`password=$(cat secret)` immediately rides into the environment of every " +
			"child process the script spawns — the `ps e`, `/proc/<pid>/environ`, and " +
			"journal of any later `| tee`, `| mail`, or `logger` call. Enabling it " +
			"script-wide to avoid a few `export` keywords leaks credentials and private " +
			"config by default. Drop the `setopt`, scope exports explicitly with " +
			"`export VAR=value`, or wrap a narrow section in `setopt LOCAL_OPTIONS; setopt " +
			"ALL_EXPORT` inside a function so the effect cannot leak past the closing " +
			"brace.",
		Check: checkZC1849,
	})
}

func checkZC1849(node ast.Node) []Violation {
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
			v := arg.String()
			if zc1849IsAllExport(v) {
				return zc1849Hit(cmd, "setopt "+v)
			}
		}
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOALLEXPORT" {
				return zc1849Hit(cmd, "unsetopt "+v)
			}
		}
	}
	return nil
}

func zc1849IsAllExport(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "ALLEXPORT"
}

func zc1849Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1849",
		Message: "`" + where + "` marks every later assignment for export — secrets " +
			"like `password=...` leak into every child's env. Drop it; use " +
			"explicit `export`, or scope inside a `LOCAL_OPTIONS` function.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

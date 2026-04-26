// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1870",
		Title:    "Warn on `setopt GLOB_ASSIGN` — RHS of `var=pattern` silently glob-expands",
		Severity: SeverityWarning,
		Description: "`GLOB_ASSIGN` is off by default in Zsh: `logs=*.log` sets `$logs` to the " +
			"literal string `*.log`, just like every other shell. Turning it on expands the " +
			"right-hand side of unquoted assignments — `logs=*.log` silently becomes the " +
			"first matching filename, `latest=backup-*` captures whatever sort-order the " +
			"filesystem returns, and any empty-match case assigns an empty string. Scripts " +
			"that port cleanly between Bash and Zsh suddenly diverge, and sensitive " +
			"assignments like `cert=~/secrets/*` can grab attacker-dropped files. Keep the " +
			"option off; use `set -A arr *.log` or explicit `arr=( *.log )` when you really " +
			"want the expansion.",
		Check: checkZC1870,
	})
}

func checkZC1870(node ast.Node) []Violation {
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
			if zc1870IsGlobAssign(arg.String()) {
				return zc1870Hit(cmd, "setopt "+arg.String())
			}
		}
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOGLOBASSIGN" {
				return zc1870Hit(cmd, "unsetopt "+v)
			}
		}
	}
	return nil
}

func zc1870IsGlobAssign(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "GLOBASSIGN"
}

func zc1870Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1870",
		Message: "`" + where + "` expands glob patterns on the RHS of `var=` — " +
			"`logs=*.log` silently captures the first match, `cert=~/secrets/*` " +
			"picks up attacker drops. Keep it off; use explicit `arr=( *.log )`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

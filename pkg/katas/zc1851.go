// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1851",
		Title:    "Warn on `unsetopt FUNCTION_ARGZERO` — `$0` inside a function stops reporting the function name",
		Severity: SeverityWarning,
		Description: "`FUNCTION_ARGZERO` is Zsh's default: inside a function or `source`d file, " +
			"`$0` holds the function/file name, which is what every `log_error \"$0: ...\"` " +
			"helper, every self-reflecting `$funcfiletrace` fallback, and every `case $0` " +
			"dispatcher expects. Turning it off reverts to POSIX-sh behaviour where `$0` " +
			"always points at the outer script — so `my_func() { echo \"${0}: bad input\" }` " +
			"silently starts logging `myscript.sh: bad input` for every function, which " +
			"makes stack-trace logs unreadable and breaks dispatchers that branch on `$0`. " +
			"Keep the option on at the script level and, if one specific helper needs the " +
			"POSIX name, reach it explicitly with `$ZSH_ARGZERO` or `$ZSH_SCRIPT`.",
		Check: checkZC1851,
	})
}

func checkZC1851(node ast.Node) []Violation {
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
			if zc1851IsFunctionArgzero(arg.String()) {
				return zc1851Hit(cmd, "unsetopt "+arg.String())
			}
		}
	case "setopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOFUNCTIONARGZERO" {
				return zc1851Hit(cmd, "setopt "+v)
			}
		}
	}
	return nil
}

func zc1851IsFunctionArgzero(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "FUNCTIONARGZERO"
}

func zc1851Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1851",
		Message: "`" + where + "` makes `$0` inside functions point at the outer " +
			"script — breaks `log \"$0: ...\"` helpers and `case $0` dispatchers. " +
			"Keep the option on; reach the script name explicitly via " +
			"`$ZSH_ARGZERO` when needed.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

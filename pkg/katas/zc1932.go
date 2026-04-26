// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1932",
		Title:    "Warn on `unsetopt GLOBAL_EXPORT` — `typeset -x` in a function stops leaking to outer scope",
		Severity: SeverityWarning,
		Description: "`GLOBAL_EXPORT` (on by default) makes `typeset -x VAR=val` inside a function " +
			"not only export `VAR` but also promote it to the outer scope, so callers and " +
			"subsequent functions see the same value. Turning it off changes the meaning of " +
			"every such assignment across the script: exports become function-local and " +
			"vanish the moment the function returns. Scripts that rely on a helper to set up " +
			"`PATH`, `VIRTUAL_ENV`, or `AWS_*` variables suddenly run commands under the old " +
			"environment. Keep the option on; if you want a temporary export, scope it with a " +
			"subshell instead of a shell-wide flip.",
		Check: checkZC1932,
	})
}

func checkZC1932(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var disabling bool
	switch ident.Value {
	case "unsetopt":
		disabling = true
	case "setopt":
		disabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1932Canonical(arg.String())
		switch v {
		case "GLOBALEXPORT":
			if disabling {
				return zc1932Hit(cmd, "unsetopt GLOBAL_EXPORT")
			}
		case "NOGLOBALEXPORT":
			if !disabling {
				return zc1932Hit(cmd, "setopt NO_GLOBAL_EXPORT")
			}
		}
	}
	return nil
}

func zc1932Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc1932Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1932",
		Message: "`" + form + "` makes `typeset -x` exports function-local — helper " +
			"functions that set `PATH`/`VIRTUAL_ENV`/`AWS_*` no longer propagate to callers. " +
			"Keep it on; scope temporary exports in a subshell instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

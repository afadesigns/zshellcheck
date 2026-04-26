// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1938",
		Title:    "Warn on `setopt POSIX_JOBS` — flips job-control semantics and `%n` scope",
		Severity: SeverityWarning,
		Description: "`POSIX_JOBS` makes Zsh's job-control spec follow POSIX: `%1` / `%n` refer " +
			"only to jobs of the current shell (forked subshells get their own job table), " +
			"`fg`/`bg` no longer accept a job ID from an outer shell, and `disown` on a " +
			"subshell's job is a no-op. Scripts that launched a background job in the parent " +
			"and then `wait %1`-ed from a `( subshell )` suddenly fail with \"no such job\". " +
			"Leave the option off in Zsh; if POSIX job semantics are required, scope them via " +
			"`emulate -LR sh` inside the single function that needs them.",
		Check: checkZC1938,
	})
}

func checkZC1938(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var enabling bool
	switch ident.Value {
	case "setopt":
		enabling = true
	case "unsetopt":
		enabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1938Canonical(arg.String())
		switch v {
		case "POSIXJOBS":
			if enabling {
				return zc1938Hit(cmd, "setopt POSIX_JOBS")
			}
		case "NOPOSIXJOBS":
			if !enabling {
				return zc1938Hit(cmd, "unsetopt NO_POSIX_JOBS")
			}
		}
	}
	return nil
}

func zc1938Canonical(s string) string {
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

func zc1938Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1938",
		Message: "`" + form + "` scopes `%n` / `fg` / `bg` / `disown` per subshell — parent " +
			"jobs become invisible inside `(…)`. Leave off; scope POSIX job semantics with " +
			"`emulate -LR sh` inside a function.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

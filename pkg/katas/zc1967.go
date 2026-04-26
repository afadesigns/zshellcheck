// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1967",
		Title:    "Warn on `setopt PROMPT_SUBST` — expansions inside `$PROMPT` evaluate command substitution every redraw",
		Severity: SeverityWarning,
		Description: "`setopt PROMPT_SUBST` turns on parameter, command, and arithmetic " +
			"substitution inside `$PS1`/`$PROMPT`/`$RPROMPT`. Any value that lands in the " +
			"prompt from an untrusted source — a git branch name, a checkout path, a " +
			"hostname in `/etc/hostname`, an env var set by a spawned tool — is reparsed " +
			"as shell code on every redraw, so a branch like `$(id>/tmp/p)` runs each time " +
			"the cursor returns. Prefer Zsh prompt escapes (`%n`, `%d`, `%~`, `%m`, " +
			"`vcs_info`) which already sanitise their inputs, or scope with `setopt " +
			"LOCAL_OPTIONS PROMPT_SUBST` inside the prompt-building function instead of " +
			"flipping the option globally.",
		Check: checkZC1967,
	})
}

func checkZC1967(node ast.Node) []Violation {
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
		v := zc1967Canonical(arg.String())
		switch v {
		case "PROMPTSUBST":
			if enabling {
				return zc1967Hit(cmd, "setopt PROMPT_SUBST")
			}
		case "NOPROMPTSUBST":
			if !enabling {
				return zc1967Hit(cmd, "unsetopt NO_PROMPT_SUBST")
			}
		}
	}
	return nil
}

func zc1967Canonical(s string) string {
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

func zc1967Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1967",
		Message: "`" + form + "` re-runs command substitution on every prompt " +
			"redraw — a branch/host/dir value with `$(…)` executes each render. Prefer " +
			"`%n`/`%d`/`%~`/`vcs_info`, or scope via `LOCAL_OPTIONS`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

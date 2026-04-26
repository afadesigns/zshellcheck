// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1928",
		Title:    "Warn on `setopt SHARE_HISTORY` — every session writes its history into every sibling session",
		Severity: SeverityWarning,
		Description: "`SHARE_HISTORY` flushes each command to `$HISTFILE` immediately and tells " +
			"all other running zsh sessions to re-read the file. A secret typed in a one-off " +
			"\"private\" terminal — `ssh user@host \"$PASS\"`, `aws sts ... --output text`, " +
			"`git push https://user:token@…` — shows up in every other terminal's `fc -l` list " +
			"seconds later. Prefer `setopt INC_APPEND_HISTORY` (append-only, per-session " +
			"isolation) and `setopt HIST_IGNORE_SPACE` so a leading space keeps the line out " +
			"of history altogether.",
		Check: checkZC1928,
	})
}

func checkZC1928(node ast.Node) []Violation {
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
		v := zc1928Canonical(arg.String())
		switch v {
		case "SHAREHISTORY":
			if enabling {
				return zc1928Hit(cmd, "setopt SHARE_HISTORY")
			}
		case "NOSHAREHISTORY":
			if !enabling {
				return zc1928Hit(cmd, "unsetopt NO_SHARE_HISTORY")
			}
		}
	}
	return nil
}

func zc1928Canonical(s string) string {
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

func zc1928Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1928",
		Message: "`" + form + "` flushes every command into every sibling zsh session — " +
			"secrets typed in one terminal surface in `fc -l` of every other. Prefer " +
			"`setopt INC_APPEND_HISTORY` plus `HIST_IGNORE_SPACE` for safer isolation.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

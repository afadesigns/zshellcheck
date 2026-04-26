// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1971",
		Title:    "Warn on `unsetopt GLOBAL_RCS` / `setopt NO_GLOBAL_RCS` — skips `/etc/zprofile`, `/etc/zshrc`, `/etc/zlogin`, `/etc/zlogout`",
		Severity: SeverityWarning,
		Description: "`GLOBAL_RCS` is on by default; only `/etc/zshenv` is sourced before it " +
			"can be toggled. Flipping the option off (either `unsetopt GLOBAL_RCS` or " +
			"`setopt NO_GLOBAL_RCS`) tells Zsh to skip `/etc/zprofile`, `/etc/zshrc`, " +
			"`/etc/zlogin`, and `/etc/zlogout` — which is exactly where admins put " +
			"corp-wide `PATH` hardening, audit hooks, umask, `HISTFILE` redirection, " +
			"and proxy variables. A login-shell script that disables the option in " +
			"`/etc/zshenv` neutralises every downstream system rc without a trace. " +
			"Keep the option on; if a specific helper needs pristine setup use " +
			"`emulate -LR zsh` inside a function or spawn `env -i zsh -f` scoped to " +
			"that helper.",
		Check: checkZC1971,
	})
}

func checkZC1971(node ast.Node) []Violation {
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
		v := zc1971Canonical(arg.String())
		switch v {
		case "GLOBALRCS":
			if !enabling {
				return zc1971Hit(cmd, "unsetopt GLOBAL_RCS")
			}
		case "NOGLOBALRCS":
			if enabling {
				return zc1971Hit(cmd, "setopt NO_GLOBAL_RCS")
			}
		}
	}
	return nil
}

func zc1971Canonical(s string) string {
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

func zc1971Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1971",
		Message: "`" + form + "` tells Zsh to skip `/etc/zprofile`, `/etc/zshrc`, " +
			"`/etc/zlogin`, `/etc/zlogout` — corp `PATH`/audit/umask/proxy config " +
			"silently dropped. Keep on; scope pristine setup with `emulate -LR zsh` " +
			"or `env -i zsh -f`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1936",
		Title:    "Warn on `setopt POSIX_ALIASES` — aliases on reserved words (`if`, `for`, …) stop expanding",
		Severity: SeverityWarning,
		Description: "Zsh by default lets `alias if='…'`, `alias function='…'`, etc. expand when " +
			"the reserved word appears in command position — the feature that makes oh-my-zsh " +
			"plugins able to hook `if` into their `preexec` chain. `setopt POSIX_ALIASES` " +
			"narrows alias expansion to plain identifiers, so any library that aliased a " +
			"reserved word silently stops being picked up. Keep the option off for " +
			"interactive Zsh; if you need POSIX parity for a specific block, wrap it with " +
			"`emulate -LR sh` instead of flipping the flag script-wide.",
		Check: checkZC1936,
	})
}

func checkZC1936(node ast.Node) []Violation {
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
		v := zc1936Canonical(arg.String())
		switch v {
		case "POSIXALIASES":
			if enabling {
				return zc1936Hit(cmd, "setopt POSIX_ALIASES")
			}
		case "NOPOSIXALIASES":
			if !enabling {
				return zc1936Hit(cmd, "unsetopt NO_POSIX_ALIASES")
			}
		}
	}
	return nil
}

func zc1936Canonical(s string) string {
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

func zc1936Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1936",
		Message: "`" + form + "` narrows alias expansion to plain identifiers — aliases " +
			"on `if`/`for`/`function` silently stop firing and any library that hooked " +
			"them breaks. Scope with `emulate -LR sh` instead of flipping globally.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

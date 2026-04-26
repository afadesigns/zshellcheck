// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1983",
		Title:    "Warn on `setopt CSH_JUNKIE_QUOTES` — single/double-quoted strings that span lines become errors",
		Severity: SeverityWarning,
		Description: "With `CSH_JUNKIE_QUOTES` off (the default), Zsh lets `\"foo\\nbar\"` and " +
			"`'line1\\nline2'` span physical lines. Setting the option on makes the " +
			"parser emit an error on the first newline inside a quoted string — which " +
			"breaks any existing multi-line SQL, JSON, or here-style payload that the " +
			"script has been inlining up to this point. Functions that are autoloaded " +
			"later or sourced from third-party helpers fail to parse, and the " +
			"diagnostic points at the closing quote, not at the option toggle. Leave " +
			"the option off; if csh-style strictness is genuinely required, scope " +
			"with `emulate -LR csh` inside the single helper that needs it.",
		Check: checkZC1983,
	})
}

func checkZC1983(node ast.Node) []Violation {
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
		v := zc1983Canonical(arg.String())
		switch v {
		case "CSHJUNKIEQUOTES":
			if enabling {
				return zc1983Hit(cmd, "setopt CSH_JUNKIE_QUOTES")
			}
		case "NOCSHJUNKIEQUOTES":
			if !enabling {
				return zc1983Hit(cmd, "unsetopt NO_CSH_JUNKIE_QUOTES")
			}
		}
	}
	return nil
}

func zc1983Canonical(s string) string {
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

func zc1983Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1983",
		Message: "`" + form + "` makes every later multi-line `\"…\"`/`'…'` an error — " +
			"inlined SQL/JSON payloads and autoloaded helpers stop parsing. Scope " +
			"csh-style strictness with `emulate -LR csh` in the one helper that " +
			"needs it.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

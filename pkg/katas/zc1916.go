// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1916",
		Title:    "Warn on `setopt NULL_GLOB` — every unmatched glob silently expands to nothing",
		Severity: SeverityWarning,
		Description: "`setopt NULL_GLOB` removes the Zsh default behaviour of erroring out when a " +
			"pattern matches nothing. Every later glob becomes silently empty instead — " +
			"`cp *.log /dest` when no `.log` files exist turns into `cp /dest` (wrong target), " +
			"`rm *.tmp` into `rm` (argv too short), and `for f in *.json` into a no-op. Reach " +
			"for the per-glob `*(N)` qualifier when you want a single pattern to tolerate a " +
			"zero match, or scope the switch with `setopt LOCAL_OPTIONS NULL_GLOB` inside the " +
			"one function that needs it.",
		Check: checkZC1916,
	})
}

func checkZC1916(node ast.Node) []Violation {
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
		v := zc1916Canonical(arg.String())
		switch v {
		case "NULLGLOB":
			if enabling {
				return zc1916Hit(cmd, "setopt NULL_GLOB")
			}
		case "NONULLGLOB":
			if !enabling {
				return zc1916Hit(cmd, "unsetopt NO_NULL_GLOB")
			}
		}
	}
	return nil
}

func zc1916Canonical(s string) string {
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

func zc1916Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1916",
		Message: "`" + form + "` makes every later unmatched glob silently empty — " +
			"`cp *.log /dest` rewrites to `cp /dest`, `rm *.tmp` becomes argv-too-short. " +
			"Use per-glob `*(N)`, or `setopt LOCAL_OPTIONS NULL_GLOB` in a function.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

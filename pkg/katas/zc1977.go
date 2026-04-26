// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1977",
		Title:    "Warn on `setopt CHASE_DOTS` — `cd ..` physically resolves before walking up, breaking logical paths",
		Severity: SeverityWarning,
		Description: "Default Zsh keeps `..` logical: from `/app/current/lib` (where " +
			"`/app/current` → `/app/releases/v5`), `cd ..` goes back to `/app/current`, " +
			"matching the user's mental model and blue/green deployment symlinks. " +
			"`setopt CHASE_DOTS` flips that — `..` first resolves the current directory " +
			"to its physical inode, so the same `cd ..` lands in `/app/releases/v5` " +
			"and the next `cd config` looks for `/app/releases/config` instead of " +
			"`/app/config`. Scripts that rely on `${PWD}` staying logical or on " +
			"`cd ../foo` matching the typed path break silently. Leave the option off; " +
			"use `cd -P` one-shot when a specific call really needs physical resolution.",
		Check: checkZC1977,
	})
}

func checkZC1977(node ast.Node) []Violation {
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
		v := zc1977Canonical(arg.String())
		switch v {
		case "CHASEDOTS":
			if enabling {
				return zc1977Hit(cmd, "setopt CHASE_DOTS")
			}
		case "NOCHASEDOTS":
			if !enabling {
				return zc1977Hit(cmd, "unsetopt NO_CHASE_DOTS")
			}
		}
	}
	return nil
}

func zc1977Canonical(s string) string {
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

func zc1977Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1977",
		Message: "`" + form + "` makes `cd ..` physically resolve before walking up — " +
			"blue/green `current` symlinks stop working for `../foo` lookups. " +
			"Keep off; use `cd -P` one-shot when physical resolution is needed.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

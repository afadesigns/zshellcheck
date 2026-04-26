// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1908",
		Title:    "Warn on `setopt MAGIC_EQUAL_SUBST` — enables tilde/param expansion on `key=value` args",
		Severity: SeverityWarning,
		Description: "`MAGIC_EQUAL_SUBST` tells Zsh that every unquoted argument of the form " +
			"`identifier=value` gets file expansion on the right-hand side, as if it were a " +
			"parameter assignment. Under the default (option off), `rsync host:dst=~/backup` " +
			"keeps the literal `~` — under the option on, the `~` expands to your home. " +
			"Flipping the option globally makes a whole class of literal CLI arguments silently " +
			"change meaning. Leave the option off; if a specific assignment truly needs " +
			"expansion, wrap it in quotes or use a temporary variable.",
		Check: checkZC1908,
	})
}

func checkZC1908(node ast.Node) []Violation {
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
		v := zc1908Canonical(arg.String())
		switch v {
		case "MAGICEQUALSUBST":
			if enabling {
				return zc1908Hit(cmd, "setopt MAGIC_EQUAL_SUBST")
			}
		case "NOMAGICEQUALSUBST":
			if !enabling {
				return zc1908Hit(cmd, "unsetopt NO_MAGIC_EQUAL_SUBST")
			}
		}
	}
	return nil
}

func zc1908Canonical(s string) string {
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

func zc1908Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1908",
		Message: "`" + form + "` gives every `key=value` argument tilde/parameter " +
			"expansion on the RHS — literal CLI args like `rsync host:dst=~/backup` " +
			"silently change. Keep it off; quote the assignment if expansion is really wanted.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

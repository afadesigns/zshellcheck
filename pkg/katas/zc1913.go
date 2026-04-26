// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1913",
		Title:    "Warn on `setopt ALIAS_FUNC_DEF` — re-enables defining functions with aliased names",
		Severity: SeverityWarning,
		Description: "Zsh's default refuses the syntax `ls () { … }` when `ls` is aliased — " +
			"because the alias expands at definition time and the function the author meant " +
			"to write never actually exists. `setopt ALIAS_FUNC_DEF` disables that guardrail: " +
			"the alias is suppressed during definition, and the function silently shadows the " +
			"alias afterwards. The combination is almost always a bug — one alias in a sourced " +
			"rc file quietly replaces the function. Keep the option off and write " +
			"`function \\ls () { … }` (quoted) if you really need to override an aliased name.",
		Check: checkZC1913,
	})
}

func checkZC1913(node ast.Node) []Violation {
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
		v := zc1913Canonical(arg.String())
		switch v {
		case "ALIASFUNCDEF":
			if enabling {
				return zc1913Hit(cmd, "setopt ALIAS_FUNC_DEF")
			}
		case "NOALIASFUNCDEF":
			if !enabling {
				return zc1913Hit(cmd, "unsetopt NO_ALIAS_FUNC_DEF")
			}
		}
	}
	return nil
}

func zc1913Canonical(s string) string {
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

func zc1913Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1913",
		Message: "`" + form + "` lets a function silently shadow an alias — one sourced " +
			"rc file replaces your function with the alias, no error surfaces. Keep it " +
			"off; quote the name if the override is intentional.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

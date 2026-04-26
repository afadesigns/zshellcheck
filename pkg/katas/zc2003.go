// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC2003",
		Title:    "Warn on `setopt KSH_ZERO_SUBSCRIPT` — `$arr[0]` stops aliasing the first element",
		Severity: SeverityWarning,
		Description: "Default Zsh treats `$arr[0]` as a quirk-compatibility alias for `$arr[1]` " +
			"— `arr=(a b c); echo $arr[0]` prints `a`, and `arr[0]=new` rewrites the " +
			"first element. `setopt KSH_ZERO_SUBSCRIPT` flips that to ksh semantics: " +
			"`$arr[0]` becomes a distinct slot (the element just before the " +
			"1-indexed head, which Zsh stores separately), so reads silently switch " +
			"to empty string and `arr[0]=new` no longer touches `$arr[1]`. Any Zsh " +
			"code that intentionally used `$arr[0]` as a shortcut breaks, and ported " +
			"Bash/ksh code that assumes 0-indexed access meets a split-world model. " +
			"Leave the option off; use `$arr[1]` explicitly when you want the first " +
			"element, and adopt `KSH_ARRAYS` scoped with `emulate -LR ksh` for " +
			"ksh-style code paths.",
		Check: checkZC2003,
	})
}

func checkZC2003(node ast.Node) []Violation {
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
		v := zc2003Canonical(arg.String())
		switch v {
		case "KSHZEROSUBSCRIPT":
			if enabling {
				return zc2003Hit(cmd, "setopt KSH_ZERO_SUBSCRIPT")
			}
		case "NOKSHZEROSUBSCRIPT":
			if !enabling {
				return zc2003Hit(cmd, "unsetopt NO_KSH_ZERO_SUBSCRIPT")
			}
		}
	}
	return nil
}

func zc2003Canonical(s string) string {
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

func zc2003Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC2003",
		Message: "`" + form + "` stops aliasing `$arr[0]` to `$arr[1]` — every later " +
			"read of `$arr[0]` silently returns empty and `arr[0]=new` stops " +
			"updating the first element. Use `$arr[1]` explicitly.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1934",
		Title:    "Warn on `setopt AUTO_NAME_DIRS` — any absolute-path parameter becomes a `~name` alias",
		Severity: SeverityWarning,
		Description: "`AUTO_NAME_DIRS` (off by default) auto-registers any parameter whose value is " +
			"an absolute directory path as a named directory — so `foo=/srv/data` immediately " +
			"makes `~foo` resolve to `/srv/data` in later expansions and in `%~` prompt " +
			"sequences. The option silently changes the meaning of `ls ~foo` across the " +
			"script and surfaces directory names in `%~` prompts that the user never opted " +
			"into. Keep the option off and call `hash -d name=/path` explicitly when a named " +
			"directory is actually wanted.",
		Check: checkZC1934,
	})
}

func checkZC1934(node ast.Node) []Violation {
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
		v := zc1934Canonical(arg.String())
		switch v {
		case "AUTONAMEDIRS":
			if enabling {
				return zc1934Hit(cmd, "setopt AUTO_NAME_DIRS")
			}
		case "NOAUTONAMEDIRS":
			if !enabling {
				return zc1934Hit(cmd, "unsetopt NO_AUTO_NAME_DIRS")
			}
		}
	}
	return nil
}

func zc1934Canonical(s string) string {
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

func zc1934Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1934",
		Message: "`" + form + "` auto-registers every absolute-path parameter as a " +
			"named dir — `foo=/srv/data` makes `~foo` expand, `%~` prompts surface names " +
			"the user never picked. Keep off; use `hash -d name=/path`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
